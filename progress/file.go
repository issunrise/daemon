package progress

import (
	"os"
)

type FileClient struct {
	file *os.File
}

func NewFileClient(name string) *FileClient {
	file, _ := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	return &FileClient{file}
}

func (c *FileClient) Update(name string) error {
	c.file.Close()
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	c.file = file
	return nil
}

func (c *FileClient) Write(b []byte) (n int, err error) {
	return c.file.Write(b)
}

func (c *FileClient) Close() error {
	return c.file.Close()
}
