package storage

type storageClient struct {
	mainPath string
}

type StorageClient interface {
	UploadImage(fileData []byte, dest, filename string) error
	DownloadFile(path string) (DownloadFileResp, error)
}

func Init(path string) StorageClient {
	if path == "" {
		path = "./etc/storage/"
	}

	return &storageClient{
		mainPath: path,
	}
}
