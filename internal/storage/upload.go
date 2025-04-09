package storage

import (
	"log"
	"os"
)

func (s *storageClient) UploadImage(fileData []byte, dest, filename string) error {
	if err := os.MkdirAll(s.mainPath+dest, os.ModePerm); err != nil {
		log.Println("[UploadImage] error while mkdir", err.Error())
		return err
	}

	if err := os.WriteFile(s.mainPath+dest+filename, fileData, 0644); err != nil {
		log.Println("[UploadImage] error while write file", err.Error())
		return err
	}

	return nil
}
