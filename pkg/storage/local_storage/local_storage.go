package localstorage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	guuid "github.com/google/uuid"
)

type LocalStorage struct {
	path string
}

func New(path string) LocalStorage {
	return LocalStorage{
		path: path,
	}
}

func (s LocalStorage) GetImage(id, quality string) ([]byte, error) {
	if quality == "" {
		quality = "100"
	}
	fileBytes, err := ioutil.ReadFile(s.path + id + "/" + quality + ".jpg")
	return fileBytes, err
}

func (s LocalStorage) SaveImage(file multipart.File, header *multipart.FileHeader) error {
	var img image.Image
	var err error
	if strings.HasSuffix(header.Filename, ".png") {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		return err
	}

	id := guuid.NewString()

	err = os.Mkdir(s.path+id, 0755)
	if err != nil {
		return err
	}

	errs := make(chan error, 4)

	for _, q := range []int{100, 75, 50, 25} {
		go func(image image.Image, quality int, id string) {
			file, err := os.Create(s.path + id + "/" + fmt.Sprint(quality) + ".jpg")
			defer file.Close()
			if err != nil {
				errs <- err
				return
			}

			err = jpeg.Encode(file, image, &jpeg.Options{Quality: quality})
			if err != nil {
				errs <- err
				return
			}
			errs <- nil
		}(img, q, id)
	}
	for i := 0; i < 4; i++ {
		err := <-errs
		if err != nil {
			return err
		}
	}

	return nil
}
