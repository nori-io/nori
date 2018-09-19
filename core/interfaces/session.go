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

package interfaces

import (
	"time"

	"context"

	"github.com/secure2work/nori/core/endpoint"
)

type SessionVerification int

const (
	NoVerify SessionVerification = iota
	WhiteList
	BlackList
)

type State string

const (
	SessionActive  State = "active"
	SessionClosed  State = "closed"
	SessionLocked  State = "locked"
	SessionBlocked State = "blocked"
	SessionExpired State = "expired"
	SessionError   State = "error"
)

const (
	SessionContextKey = "nori.session.id"
)

type Session interface {
	Get(key []byte, data interface{}) error
	Save(key []byte, data interface{}, exp time.Duration) error
	Delete(key []byte) error
	SessionId(ctx context.Context) []byte
	Verify(verify SessionVerification) endpoint.Middleware
}
