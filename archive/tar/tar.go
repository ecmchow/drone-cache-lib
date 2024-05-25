package tar

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ecmchow/drone-cache-lib/archive"
	log "github.com/sirupsen/logrus"
)

type tarArchive struct{}

// New creates an archive that uses the .tar file format.
func New() archive.Archive {
	return &tarArchive{}
}

func (a *tarArchive) Pack(srcs []string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	// Loop through each source
	var fwErr error
	for _, s := range srcs {
		// ensure the src actually exists before trying to tar it
		if _, err := os.Stat(s); err != nil {
			return err
		}

		// walk path
		fwErr = filepath.Walk(s, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}

			var link string
			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
				log.Debugf("Symbolic link found at %s to %s", path, link)

				// Rewrite header for SymLink
				header, err = tar.FileInfoHeader(fi, link)
				if err != nil {
					return err
				}
			}

			header.Name = strings.TrimPrefix(filepath.ToSlash(path), "/")

			if err = tw.WriteHeader(header); err != nil {
				return err
			}

			if !fi.Mode().IsRegular() {
				log.Debugf("Directory found at %s", path)
				return nil
			}

			log.Debugf("File found at %s", path)

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		})

		if fwErr != nil {
			return fwErr
		}
	}

	return fwErr
}

func (a *tarArchive) Unpack(dst string, r io.Reader) error {
	log.Println("Unpack started")
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			log.Println("No more files found, returning")
			return nil
		case err != nil:
			log.Printf("Error reading next header: %v", err)
			return err
		case header == nil:
			log.Println("Header is nil, skipping")
			continue
		}

		target := filepath.Join(dst, header.Name)
		log.Printf("Processing file: %s", target)

		switch header.Typeflag {
		case tar.TypeSymlink:
			log.Printf("Symlink found at %s", target)
			_, err := os.Stat(target)
			if err == nil {
				log.Printf("Failed to create symlink because file already exists at %s", target)
				return fmt.Errorf("Failed to create symlink because file already exists at %s", target)
			}

			log.Printf("Creating link %s to %s", target, header.Linkname)
			err = os.Symlink(header.Linkname, target)
			if err != nil {
				log.Printf("Failed creating link %s to %s", target, header.Linkname)
				return err
			}

		case tar.TypeDir:
			log.Printf("Directory found at %s", target)
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			log.Printf("File found at %s", target)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(f, tr)
			f.Close()
			if err != nil {
				log.Printf("Error copying file contents: %v", err)
				return err
			}

			err = os.Chtimes(target, time.Now(), header.ModTime)
			if err != nil {
				log.Printf("Error changing file times: %v", err)
				return err
			}
		}
	}
}
