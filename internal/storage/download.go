package storage

import (
	"log"
	"os"
)

type DownloadFileResp struct {
	File *os.File
}

func (s *client) DownloadFile(path string) (DownloadFileResp, error) {
	f, err := os.Open(s.mainPath + path)
	if err != nil {
		log.Println("[DownloadFile] image not found", err.Error())
		return DownloadFileResp{}, err
	}

	return DownloadFileResp{
		File: f,
	}, nil
}
