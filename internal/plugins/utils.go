/*
Copyright 2019-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugins

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
)

func LogFieldsMeta(m meta.Meta) []logger.Field {
	return []logger.Field{
		{"plugin", string(m.Id().String())},
		{"interface", m.GetInterface().String()},
	}
}
