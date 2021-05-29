package file

import "github.com/nori-io/nori/internal/domain/repository"

func New() repository.FileRepository {
	return &FileRepository{}
}
