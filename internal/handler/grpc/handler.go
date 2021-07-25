package grpc

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	pkgmeta "github.com/nori-io/common/v5/pkg/meta"
	"github.com/nori-io/nori-grpc/pkg/api/proto"
	"github.com/nori-io/nori/internal/domain/enum"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	errors2 "github.com/nori-io/nori/pkg/nori/domain/errors"
	"go.uber.org/dig"
)

type Handler struct {
	proto.UnimplementedNoriServer
	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginService    service.PluginService
	PluginManager    service.PluginManager
	View             View
}

type HandlerParams struct {
	dig.In

	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginManager    service.PluginManager
	PluginService    service.PluginService
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		FileService:      params.FileService,
		InstalledService: params.InstalledService,
		PluginManager:    params.PluginManager,
		PluginService:    params.PluginService,
		View:             View{},
	}
}

//
func (h Handler) PluginEnable(ctx context.Context, in *proto.PluginRequest) (*proto.Reply, error) {
	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	if err := h.PluginManager.Enable(ctx, id); err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	return &proto.Reply{}, nil
}

//
func (h Handler) PluginDisable(ctx context.Context, in *proto.PluginRequest) (*proto.Reply, error) {
	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	if err := h.PluginManager.Disable(ctx, id); err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	return &proto.Reply{}, nil
}

//
func (h Handler) PluginGet(ctx context.Context, in *proto.PluginRequest) (*proto.Reply, error) {
	return nil, nil
}

//
func (h Handler) PluginInstall(ctx context.Context, in *proto.PluginInstallRequest) (*proto.Reply, error) {
	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	if err := h.PluginManager.Install(ctx, id); err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	return &proto.Reply{}, nil
}

//
func (h Handler) PluginInterface(ctx context.Context, in *proto.PluginInterfaceRequest) (*proto.PluginListReply, error) {
	list := []*entity.Plugin{}
	//for _, p := range h.PluginManager.GetAll() {
	//	if p.Meta().GetInterface().String() == in.Interface || p.Meta().GetInterface().Name() == in.Interface {
	//		list = append(list, p)
	//	}
	//}

	plugins := []*proto.Plugin{}
	for _, p := range list {
		plugins = append(plugins, h.View.Plugin(p))
	}

	return &proto.PluginListReply{
		Plugin: plugins,
		Error:  nil,
	}, nil
}

//
func (h Handler) PluginList(ctx context.Context, in *proto.PluginListRequest) (*proto.PluginListReply, error) {
	items := []*entity.Plugin{}

	if in.FlagInstallable {
		// todo ???
	}

	if in.FlagInstalled {
		installed, err := h.InstalledService.GetAll()
		if err != nil {
			return &proto.PluginListReply{
				Error: &proto.Error{
					Code:    "500",
					Message: err.Error(),
				},
			}, nil
		}

		plugins := []*proto.Plugin{}
		for _, i := range installed {
			plugins = append(plugins, &proto.Plugin{
				Meta: &proto.Meta{
					Id: &proto.ID{
						PluginId: string(i.ID.GetID()),
						Version:  i.ID.GetVersion(),
					},
				},
				IsActive:      false,
				IsInstallable: false,
				Status:        0,
			})
		}
	}

	if in.FlagError {
		// todo ???
	}

	if in.FlagRunning {
		running, err := h.PluginManager.GetByFilter(service.GetByFilterData{State: enum.Running})
		if err != nil {
			return nil, err
		}
		items = append(items, running...)
	}

	if in.FlagStopped {
		stopped, err := h.PluginManager.GetByFilter(service.GetByFilterData{State: enum.Stopped})
		if err != nil {
			return nil, err
		}
		items = append(items, stopped...)
	}

	plugins := []*proto.Plugin{}
	for _, item := range items {
		plugins = append(plugins, h.View.Plugin(item))
	}

	return &proto.PluginListReply{
		Plugin: plugins,
		Error:  nil,
	}, nil
}

//
func (h Handler) PluginMeta(ctx context.Context, in *proto.PluginMetaRequest) (*proto.PluginMetaReply, error) {
	id, err := h.toMetaID(in.Id)
	if err != nil {
		return nil, err
	}

	plugin, err := h.PluginService.Get(id)
	if err != nil {
		return nil, err
	}

	return &proto.PluginMetaReply{
		Meta: h.View.Meta(plugin.Meta()),
	}, nil
}

// todo
func (h Handler) PluginPull(ctx context.Context, in *proto.PluginPullRequest) (*proto.Reply, error) {
	return nil, nil
}

// todo
func (h Handler) PluginRemove(ctx context.Context, in *proto.PluginRequest) (*proto.Reply, error) {
	return &proto.Reply{
		Error: &proto.Error{
			Code:    "500",
			Message: "not implemented",
		},
	}, nil
}

//
func (h Handler) PluginStart(ctx context.Context, in *proto.PluginStartRequest) (*proto.Reply, error) {
	if in.FlagAll {
		err := h.PluginManager.StartAll(ctx)
		if err != nil {
			return &proto.Reply{Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			}}, nil
		}
		return &proto.Reply{}, nil
	}

	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}

	err = h.PluginManager.Start(ctx, id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}

	return &proto.Reply{}, nil
}

//
func (h Handler) PluginStop(ctx context.Context, in *proto.PluginStopRequest) (*proto.Reply, error) {
	if in.FlagAll {
		err := h.PluginManager.StopAll(ctx)
		if err != nil {
			return &proto.Reply{
				Error: &proto.Error{
					Code:    "500",
					Message: err.Error(),
				},
			}, nil
		}
		return &proto.Reply{}, nil
	}

	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}

	err = h.PluginManager.Stop(ctx, id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}

	return &proto.Reply{}, nil
}

//
func (h Handler) PluginUninstall(ctx context.Context, in *proto.PluginUninstallRequest) (*proto.Reply, error) {
	id, err := h.toMetaID(in.Id)
	if err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}

	if err := h.PluginManager.UnInstall(ctx, id); err != nil {
		return &proto.Reply{
			Error: &proto.Error{
				Code:    "500",
				Message: err.Error(),
			},
		}, nil
	}
	return &proto.Reply{}, nil
}

// todo
func (h Handler) PluginUpload(stream proto.Nori_PluginUploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	//@todo check if file already exists (os.Stat, isDir)

	pluginData := bytes.Buffer{}
	pluginSize := 0

	for {
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Reply{})
		}
		if err != nil {
			return err
		}

		chunk := req.GetChunk()
		size := len(chunk)

		log.Println("received a chunk with size:", size)

		pluginSize += size
		//@todo code down in the comment
		/*	if pluginSize > maxPluginSize {
			log.Println(status.Errorf(codes.InvalidArgument, "plugin is too large: %d > %d", pluginSize, maxPluginSize))
			return err
		}*/

		_, err = pluginData.Write(chunk)
		if err != nil {
			return err
		}
	}

	_, err = h.FileService.Create(req.GetName(), pluginData)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&proto.Reply{})

}

// todo
func (h Handler) PluginDownload(*proto.PluginDownloadRequest, proto.Nori_PluginDownloadServer) error {
	return nil
}

//config
// todo
func (h Handler) ConfigGet(ctx context.Context, in *proto.ConfigGetRequest) (*proto.ConfigGetReply, error) {
	//id, err := h.toMetaID(in.Id)
	//if err != nil {
	//	return nil, err
	//}

	//vars := h.ConfigManager.PluginVariables(id)

	return &proto.ConfigGetReply{}, nil
}

// todo
func (h Handler) ConfigSet(ctx context.Context, in *proto.ConfigSetRequest) (*proto.ConfigSetReply, error) {
	return nil, nil
}

// todo
func (h Handler) ConfigUpload(ctx context.Context, in *proto.ConfigUploadRequest) (*proto.ConfigUploadReply, error) {
	return nil, nil
}

func (h Handler) toMetaID(id *proto.ID) (meta.ID, error) {
	if id == nil {
		return nil, errors2.NotFound{}
	}

	return pkgmeta.ID{
		ID:      meta.PluginID(id.PluginId),
		Version: id.Version,
	}, nil
}
