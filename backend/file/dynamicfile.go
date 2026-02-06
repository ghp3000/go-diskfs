package file

import (
	"errors"
	"fmt"
	"github.com/ghp3000/go-diskfs/backend"
	"io/fs"
	"os"
)

const expansionStep int64 = 10 * 1024 * 1024 // 100MB
type DynamicFile struct {
	storage *os.File
}

func CreateDynamicStorage(pathName string) (backend.Storage, error) {
	if pathName == "" {
		return nil, errors.New("must pass device name")
	}
	f, err := os.OpenFile(pathName, os.O_RDWR|os.O_EXCL|os.O_CREATE, 0o666)
	if err != nil {
		return nil, fmt.Errorf("could not create device %s: %w", pathName, err)
	}
	return New(NewDynamicFile(f), false), nil
}

func NewDynamicFile(f *os.File) fs.File {
	return &DynamicFile{
		storage: f,
	}
}
func (f *DynamicFile) Stat() (fs.FileInfo, error) {
	return f.storage.Stat()
}
func (f *DynamicFile) Read(b []byte) (int, error) {
	return f.storage.Read(b)
}
func (f *DynamicFile) Close() error {
	return f.storage.Close()
}
func (f *DynamicFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.storage.ReadAt(p, off)
}
func (f *DynamicFile) Seek(offset int64, whence int) (int64, error) {
	return f.storage.Seek(offset, whence)
}
func (f *DynamicFile) WriteAt(p []byte, off int64) (n int, err error) {
	stat, err := f.storage.Stat()
	if err != nil {
		return 0, err
	}
	required := off + int64(len(p))
	if stat.Size() < required {
		newSize := required + expansionStep
		if err := f.storage.Truncate(newSize); err != nil {
			return 0, err
		}
	}
	return f.storage.WriteAt(p, off)
}
