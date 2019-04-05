// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package grpc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/nori-io/nori-common/meta"

	"github.com/nori-io/nori/core/plugins"
	"github.com/nori-io/nori/proto"
)

const (
	MaxMessageSize int = 100 * 1024 * 1024
)

var (
	NotSecure = errors.New("Need safe gRPC connect")
)

type Server struct {
	pluginDirs  []string
	gRPCAddress string
	gRPCEnable  bool

	certFile string
	keyFile  string

	pluginManager plugins.Manager
	passkey       *Passkey
	grpcServer    *grpc.Server
	gShutdown     chan struct{}
	secure        bool
	done          bool
	log           *logrus.Logger
}

func NewServer(
	dirs []string,
	addr string,
	enable bool,
	pluginManager plugins.Manager,
	log *logrus.Logger,
) *Server {
	return &Server{
		pluginManager: pluginManager,
		pluginDirs:    dirs,
		gRPCAddress:   addr,
		gRPCEnable:    enable,
		gShutdown:     make(chan struct{}, 1),
		log:           log,
	}
}

func (s *Server) SetCertificates(cert, key string) {
	s.certFile = cert
	s.keyFile = key
}

func (s *Server) Run() error {
	var err error

	s.passkey, err = NewPasskey()
	if err != nil {
		return err
	}

	s.log.Infof("Passkey: %s", s.passkey)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(s *Server) {
		defer wg.Done()
		for !s.done {
			listener, err := net.Listen("tcp", s.gRPCAddress)
			if err != nil {
				s.log.Fatal(err)
			}

			var opts []grpc.ServerOption

			opts = append(opts, grpc.MaxMsgSize(MaxMessageSize))

			if opt, err := s.CheckTLS(); err == nil {
				opts = append(opts, opt)
				s.secure = true
			}

			s.grpcServer = grpc.NewServer(opts...)
			commands.RegisterCommandsServer(s.grpcServer, s)
			logrus.WithField("Secure", s.secure).Infof("Starting gRPC server on %s", s.gRPCAddress)
			s.grpcServer.Serve(listener)
		}
	}(s)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func(s *Server) {
		for {
			select {
			case sig := <-signalCh:
				logrus.Infof("Graceful stop gRPC server with signal: %s", sig)
				s.Stop()
			case <-s.gShutdown:
				s.grpcServer.GracefulStop()
			}
		}
	}(s)

	wg.Wait()

	return nil
}

func (s *Server) Stop() {
	s.done = true
	s.gShutdown <- struct{}{}
}

func (s Server) GetPasskey() string {
	return s.passkey.String()
}

func (s Server) GetSecure() bool {
	return s.secure
}

func (s Server) PluginListCommand(_ context.Context, _ *commands.PluginListRequest) (*commands.PluginListReply, error) {
	if !s.secure {
		return nil, NotSecure
	}
	reply := new(commands.PluginListReply)
	reply.Data = make([]*commands.PluginList, 0)

	for _, m := range s.pluginManager.Metas(plugins.FilterRunnable) {
		reply.Data = append(reply.Data, &commands.PluginList{
			Id:     m.Id().String(),
			Name:   m.GetDescription().Name,
			Author: m.GetAuthor().Name,
		})
	}
	return reply, nil
}

func (s Server) PluginGetCommand(_ context.Context, c *commands.PluginGetRequest) (*commands.ErrorReply, error) {
	if !s.secure {
		return nil, NotSecure
	}
	toolchain, err := SetupToolChain()
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	toolchain.InstallDependencies = c.GetInstallDependencies()
	toolchain.PluginDir = s.pluginDirs[0]
	err = toolchain.Do(c.GetUri())
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	} else {
		return &commands.ErrorReply{
			Status: true,
			Error:  "",
		}, nil
	}
}

func (s Server) PluginRemoveCommand(_ context.Context, c *commands.PluginRemoveRequest) (*commands.ErrorReply, error) {
	// FIXME: implement
	return nil, nil
}

func (s Server) PluginMetaCommand(_ context.Context, req *commands.PluginMetaRequest) (*commands.PluginMetaReply, error) {
	return nil, nil
}

func (s Server) PluginInstallCommand(ctx context.Context, c *commands.PluginInstallRequest) (*commands.ErrorReply, error) {
	parts := strings.Split(c.Id, ":")
	if len(parts) != 2 {
		err := fmt.Errorf("ID does not contain version information")
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	id := meta.ID{
		ID:      meta.PluginID(parts[0]),
		Version: parts[1],
	}
	err := s.pluginManager.Install(id, ctx)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	return &commands.ErrorReply{
		Status: true,
		Error:  "",
	}, nil
}

func (s Server) PluginUninstallCommand(ctx context.Context, c *commands.PluginUninstallRequest) (*commands.ErrorReply, error) {
	parts := strings.Split(c.Id, ":")
	if len(parts) != 2 {
		err := fmt.Errorf("ID does not contain version information")
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	id := meta.ID{
		ID:      meta.PluginID(parts[0]),
		Version: parts[1],
	}
	err := s.pluginManager.UnInstall(id, ctx)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	return &commands.ErrorReply{
		Status: true,
		Error:  "",
	}, nil
}

func (s Server) PluginUploadCommand(_ context.Context, c *commands.PluginUploadRequest) (*commands.ErrorReply, error) {
	path := filepath.Join(s.pluginDirs[0], c.Name)
	if fileExists(path) {
		s.log.Info("File exist, overwrites")
	}

	err := os.MkdirAll(s.pluginDirs[0], os.ModePerm)
	if err != nil {
		s.log.Errorf("can't upload plugin %s: %s", c.Name, err.Error())
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	f, err := os.Create(path)
	if err != nil {
		s.log.Errorf("can't upload plugin %s: %s", c.Name, err.Error())
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	f.Write(c.So)
	f.Close()

	s.log.Infof("plugin %s uploaded", c.Name)

	if m, _ := s.pluginManager.AddFile(path); m != nil {
		s.log.Infof(
			"Found: '%s' by '%s'",
			m.Id().String(),
			m.GetAuthor(),
		)
	} else {
		s.log.Errorf("can't load plugin %s", c.Name)
	}

	return &commands.ErrorReply{
		Status: true,
		Error:  "",
	}, nil
}

func (s Server) CertsUploadCommand(_ context.Context, c *commands.CertsUploadRequest) (*commands.ErrorReply, error) {
	size := int(c.Key[:1][0])
	hmac := c.Key[1 : size+1]
	c.Key = c.Key[size+1:]
	keyBody, err := s.passkey.Decrypt(c.Key, hmac)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	size = int(c.Pem[:1][0])
	hmac = c.Pem[1 : size+1]
	c.Pem = c.Pem[size+1:]
	certBody, err := s.passkey.Decrypt(c.Pem, hmac)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	err = os.MkdirAll(filepath.Dir(s.keyFile), os.ModePerm)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	err = os.MkdirAll(filepath.Dir(s.certFile), os.ModePerm)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	fKey, err := os.Create(s.keyFile)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	defer fKey.Close()

	_, err = fKey.Write(keyBody)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	fCert, err := os.Create(s.certFile)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}
	defer fCert.Close()

	_, err = fCert.Write(certBody)
	if err != nil {
		return &commands.ErrorReply{
			Status: false,
			Error:  err.Error(),
		}, err
	}

	s.gShutdown <- struct{}{}

	return &commands.ErrorReply{
		Status: true,
		Error:  "",
	}, nil
}

func (s Server) SendPingCommand(_ context.Context, ping *commands.PingRequest) (*commands.PongReply, error) {
	return &commands.PongReply{Message: ping.Message}, nil
}

func (s Server) CheckTLS() (grpc.ServerOption, error) {
	if len(s.certFile) > 0 && len(s.keyFile) > 0 &&
		fileExists(s.certFile) && fileExists(s.keyFile) {
		creds, err := credentials.NewServerTLSFromFile(s.certFile, s.keyFile)
		if err != nil {
			return nil, err
		} else {
			return grpc.Creds(creds), nil
		}
	}
	return nil, errors.New("Bad certs")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
