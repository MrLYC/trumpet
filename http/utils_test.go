package http_test

import (
	"github.com/spf13/afero"

	"github.com/mrlyc/trumpet/http"
)

var rawFS afero.Fs

func writeFile(name string, content []byte) error {
	file, err := http.FS.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}

func setupFs(fs afero.Fs) {
	rawFS = http.FS
	http.FS = fs
}

func tearDownFs() {
	http.FS = rawFS
}
