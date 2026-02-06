package main

import (
	"fmt"
	"github.com/ghp3000/go-diskfs/backend/file"
	"github.com/ghp3000/go-diskfs/filesystem/squashfs"
	"os"
)

func main() {
	f, err := os.Create("D:\\temp.sfs")
	if err != nil {
		panic(err)
	}
	storage := file.New(file.NewDynamicFile(f), false)

	var diskSize int64 = 10 * 1024 * 1024 // 10 MB
	fs, err := squashfs.Create(storage, diskSize, 0, 131072)
	if err != nil {
		panic(err)
	}
	//err = fs.AddFile("D:\\TDDOWNLOAD\\hello\\go.mod", "go.mod")
	err = fs.AddDir("D:\\TDDOWNLOAD\\hello", "")
	//rw, err := fs.OpenFile("demo.txt", os.O_CREATE|os.O_RDWR)
	//if err != nil {
	//	panic(err)
	//}
	//content := []byte("demo")
	//_, err = rw.Write(content)
	//rw.Close()

	fs.Finalize(squashfs.FinalizeOptions{Compression: squashfs.NewCompressorZstd(15)})
	fs.Close()
	size := fs.Size()
	fmt.Println("Finalized size:", size)
	f.Truncate(int64(size))
	f.Close()
}
