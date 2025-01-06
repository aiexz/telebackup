package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log/slog"
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
			slog.Error("Some error occured", "err", err)
			// IDK what are these errors so we just print them
			return err
		}

		if header.Mode&0400 == 0 {
			slog.Debug("Skipping file", "file", file, "reason", "not readable rights")
			return nil
		}

		relPath, err := filepath.Rel(baseDir, file)
		if err != nil {
			slog.Error("Some error occured", "err", err)
			return err
		}

		header.Name = relPath

		// write header
		if err := tw.WriteHeader(header); err != nil {
			slog.Error("Some error occured", "err", err)
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				slog.Error("Some error occured", "err", err)
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				slog.Error("Some error occured", "err", err)
				return err
			}
		}
		return nil
	})
	if err := tw.Close(); err != nil {
		slog.Error("Some error occured", "err", err)
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		slog.Error("Some error occured", "err", err)
		return err
	}
	return nil

}
