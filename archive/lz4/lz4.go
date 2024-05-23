package lz4

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"io"

	"github.com/ecmchow/drone-cache-lib/archive"
	"github.com/ecmchow/drone-cache-lib/archive/tar"
	"github.com/pierrec/lz4/v4"
)

type lz4Archive struct{}

// New creates an archive that uses the .tar.gz file format.
func New() archive.Archive {
	return &lz4Archive{}
}

func (a *lz4Archive) Pack(srcs []string, w io.Writer) error {
	zw := lz4.NewWriter(w)
	defer zw.Close()

	taP := tar.New()

	err := taP.Pack(srcs, zw)

	return err
}

func (a *lz4Archive) Unpack(dst string, r io.Reader) error {
	zr := lz4.NewReader(r)

	taU := tar.New()

	fwErr := taU.Unpack(dst, zr)

	return fwErr
}
