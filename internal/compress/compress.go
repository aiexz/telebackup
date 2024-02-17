package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// CompressPath compresses the given path to the given writer using tar and gzip
func CompressPath(targetPath string, buf io.Writer) error {
	// check if targetPath is directory
	_, err := os.Stat(targetPath)
	if err != nil {
		return err
	}
	//get parent directory of targetPath
	baseDir := path.Dir(targetPath)
	// taken from https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)
	_ = filepath.Walk(targetPath, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if header.Mode&0400 == 0 {
			fmt.Println("Skipping file", file, "as it is not readable")
			return nil
		}

		relPath, err := filepath.Rel(baseDir, file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		header.Name = relPath

		// write header
		if err := tw.WriteHeader(header); err != nil {
			fmt.Println(err)
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				fmt.Println(err)
				return err
			}
		}
		return nil
	})
	if err := tw.Close(); err != nil {
		fmt.Println(err)
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}
