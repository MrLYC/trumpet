package http

import (
	"github.com/spf13/afero"
)

// FS :
var FS afero.Fs

func init() {
	FS = afero.NewOsFs()
}
