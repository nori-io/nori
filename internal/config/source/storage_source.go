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

package source

import "github.com/nori-io/nori-common/v2/storage"

type src struct {
	s storage.Storage
}

func (s *src) Get(key string) interface{} {

}

func (s *src) IsSet(key string) bool {

}

// fields
func (s *src) Bool(key string) (bool, error) {

}

func (s *src) Float(key string) (float64, error) {

}

func (s *src) Int(key string) (int, error) {

}

func (s *src) Int8(key string) (int8, error) {

}

func (s *src) Int32(key string) (int32, error) {

}

func (s *src) Int64(key string) (int64, error) {

}

func (s *src) Slice(key string) ([]interface{}, error) {

}

func (s *src) SliceInt(key string) ([]int, error) {

}

func (s *src) SliceString(key string) ([]string, error) {

}

func (s *src) String(key string) (string, error) {

}

func (s *src) StringMap(key string) map[string]interface{} {

}

func (s *src) StringMapInt(key string) map[string]int {

}

func (s *src) StringMapSliceString(key string) map[string][]string {

}

func (s *src) StringMapString(key string) map[string]string {

}

func (s *src) UInt(key string) (uint, error) {

}

func (s *src) UInt32(key string) (uint32, error) {

}

func (s *src) UInt64(key string) (uint64, error) {

}
