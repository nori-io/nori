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

package types

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/pkg/errors"
)

type MetaList []meta.Meta

func (ml MetaList) Find(id meta.ID) (meta.Meta, error) {
	for _, m := range ml {
		if m.Id() == id {
			return m, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (ml *MetaList) Add(m meta.Meta) {
	*ml = append(*ml, m)
}

func (ml *MetaList) Remove(id meta.ID) {
	for i, v := range *ml {
		if v.Id() == id {
			*ml = append((*ml)[:i], (*ml)[i+1:]...)
		}
	}
}
