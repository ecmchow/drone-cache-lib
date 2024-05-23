package zstd

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"io"

	"github.com/klauspost/compress/zstd"

	"github.com/ecmchow/drone-cache-lib/archive"
	"github.com/ecmchow/drone-cache-lib/archive/tar"
)

type zstdArchive struct{}

// New creates an archive that uses the .tar.gz file format.
func New() archive.Archive {
	return &zstdArchive{}
}

func (a *zstdArchive) Pack(srcs []string, w io.Writer) error {
	zw, err := zstd.NewWriter(w)

	if err != nil {
		return err
	}

	defer zw.Close()

	taP := tar.New()

	pErr := taP.Pack(srcs, zw)

	return pErr
}

func (a *zstdArchive) Unpack(dst string, r io.Reader) error {
	zr, err := zstd.NewReader(r)

	if err != nil {
		return err
	}

	taU := tar.New()

	fwErr := taU.Unpack(dst, zr)

	return fwErr
}
