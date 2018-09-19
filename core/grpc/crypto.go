// Copyright © 2018 Secure2Work info@secure2work.com
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
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

const (
	PasskeyLength int = 32
)

type Passkey struct {
	bs    []byte
	block cipher.Block
}

func NewPasskey() (*Passkey, error) {
	pk := new(Passkey)

	pk.bs = make([]byte, PasskeyLength*2)
	// first half for aes, second half for hmac
	_, err := rand.Read(pk.bs)
	if err != nil {
		return nil, err
	}

	pk.block, err = aes.NewCipher(pk.bs[:PasskeyLength])
	if err != nil {
		return nil, err
	}

	return pk, nil
}

func PasskeyFromString(pk string) (*Passkey, error) {
	if len(pk) != PasskeyLength*4 {
		return nil, errors.New("Bad passkey length")
	}

	bs, err := hex.DecodeString(pk)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(bs[:PasskeyLength])
	if err != nil {
		return nil, err
	}

	return &Passkey{bs: bs, block: block}, nil
}

func (pk Passkey) String() string {
	return hex.EncodeToString(pk.bs)
}

func (pk Passkey) Encrypt(bs []byte) ([]byte, []byte, error) {
	ciphertext := make([]byte, aes.BlockSize+len(bs))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCFBEncrypter(pk.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], bs)

	return ciphertext, pk.hmac(ciphertext), nil
}

func (pk Passkey) Decrypt(ciphertext []byte, hmactext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("Ciphertext too short")
	}

	if !hmac.Equal(pk.hmac(ciphertext), hmactext) {
		return nil, errors.New("malformed signature")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(pk.block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func (pk Passkey) hmac(bs []byte) []byte {
	h := hmac.New(sha256.New, pk.bs[PasskeyLength:])
	h.Write(bs)
	return h.Sum(nil)
}
