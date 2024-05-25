package s2

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"io"

	"github.com/klauspost/compress/s2"

	"github.com/ecmchow/drone-cache-lib/archive"
	"github.com/ecmchow/drone-cache-lib/archive/tar"
)

type s2Archive struct{}

// New creates an archive that uses the .tar.gz file format.
func New() archive.Archive {
	return &s2Archive{}
}

func (a *s2Archive) Pack(srcs []string, w io.Writer) error {
	zw := s2.NewWriter(w)
	defer zw.Close()

	taP := tar.New()

	pErr := taP.Pack(srcs, zw)

	return pErr
}

func (a *s2Archive) Unpack(dst string, r io.Reader) error {
	zr := s2.NewReader(r)

	taU := tar.New()

	fwErr := taU.Unpack(dst, zr)

	return fwErr
}
