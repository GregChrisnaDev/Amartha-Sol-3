package storage

import "os"

type client struct {
	mainPath string
}

type Client interface {
	UploadImage(fileData []byte, dest, filename string) error
	DownloadFile(path string) (DownloadFileResp, error)
	GetMainPath() string
}

func Init() Client {
	path := os.Getenv("DEFAULT_STORAGE")
	if path == "" {
		path = "./etc/storage/"
	}

	return &client{
		mainPath: path,
	}
}

func (c *client) GetMainPath() string {
	return c.mainPath
}
