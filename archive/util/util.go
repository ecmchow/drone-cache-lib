package util

import (
	"fmt"
	"strings"

	"github.com/ecmchow/drone-cache-lib/archive"
	"github.com/ecmchow/drone-cache-lib/archive/lz4"
	"github.com/ecmchow/drone-cache-lib/archive/tar"
	"github.com/ecmchow/drone-cache-lib/archive/tgz"
	"github.com/ecmchow/drone-cache-lib/archive/zstd"
)

// FromFilename determines the archive format to use based on the name.
func FromFilename(name string) (archive.Archive, error) {
	if strings.HasSuffix(name, ".tar") {
		return tar.New(), nil
	}

	if strings.HasSuffix(name, ".tgz") || strings.HasSuffix(name, ".tar.gz") {
		return tgz.New(), nil
	}

	if strings.HasSuffix(name, ".lz4") {
		return lz4.New(), nil
	}

	if strings.HasSuffix(name, ".zst") {
		return zstd.New(), nil
	}

	return nil, fmt.Errorf("Unknown file format for archive %s", name)
}
