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
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFromString(t *testing.T) {
	assert := assert.New(t)

	pk, err := NewPasskey()
	assert.NotNil(pk)
	assert.Nil(err)

	pkString := pk.String()

	fromStr, err := PasskeyFromString(pkString)
	assert.NotNil(pk)
	assert.Nil(err)

	assert.Equal(pk.bs, fromStr.bs)
}

func TestCrypto(t *testing.T) {
	assert := assert.New(t)

	pk, err := NewPasskey()
	assert.NotNil(pk)
	assert.Nil(err)

	testData := make([]byte, 512)
	_, err = rand.Read(pk.bs)
	assert.Nil(err)

	ciphertext, hash, err := pk.Encrypt(testData)
	assert.NotEmpty(ciphertext)
	assert.NotEmpty(hash)
	assert.Nil(err)

	data, err := pk.Decrypt(ciphertext, hash)
	assert.NotEmpty(data)
	assert.Nil(err)

	assert.Equal(testData, data)

	ciphertext[256] = ciphertext[256] * 2

	data, err = pk.Decrypt(ciphertext, hash)
	assert.Nil(data)
	assert.Equal(err.Error(), "malformed signature")
}
