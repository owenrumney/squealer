package tests

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func unpackTestRepo(tarpath string) (string, error) {
	target, err := ioutil.TempDir(os.TempDir(), "squealer_*")
	if err != nil {
		return "", err
	}

	reader, err := os.Open(tarpath)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return "", err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return "", err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return "", err
		}
	}
	return target, nil
}
