package storage

type client struct {
	mainPath string
}

type Client interface {
	UploadImage(fileData []byte, dest, filename string) error
	DownloadFile(path string) (DownloadFileResp, error)
	GetMainPath() string
}

func Init(path string) Client {
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
