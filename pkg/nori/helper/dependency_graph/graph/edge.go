/*
Copyright 2018-2020 The Nori Authors.
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

package graph

import "github.com/nori-io/common/v5/pkg/domain/meta"

type Edge interface {
	From() meta.ID
	To() meta.ID
}

type edge struct {
	from meta.ID
	to   meta.ID
}

func (e *edge) From() meta.ID {
	return e.from
}

func (e *edge) To() meta.ID {
	return e.to
}
