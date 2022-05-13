package storage

import "mime/multipart"

type Storage interface {
	SaveImage(file multipart.File, header *multipart.FileHeader) error
	GetImage(id, quality string) ([]byte, error)
}
