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
