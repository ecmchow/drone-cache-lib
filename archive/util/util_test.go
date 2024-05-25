package util

import (
	"testing"

	"github.com/franela/goblin"
)

func TestArchive(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("FromFilename", func() {
		g.It("Should return tarArchive for .tar", func() {
			_, err := FromFilename("filename.tar")
			g.Assert(err == nil).IsTrue("failed to determine .tar suffix")
		})

		g.It("Should return tgzArchive for .tgz", func() {
			_, err := FromFilename("filename.tgz")
			g.Assert(err == nil).IsTrue("failed to determine .tgz suffix")
		})

		g.It("Should return tgzArchive for .tar.gz", func() {
			_, err := FromFilename("filename.tar.gz")
			g.Assert(err == nil).IsTrue("failed to determine .tar.gz suffix")
		})

		g.It("Should return s2Archive for .tsz", func() {
			_, err := FromFilename("filename.tsz")
			g.Assert(err == nil).IsTrue("failed to determine .tsz suffix")
		})

		g.It("Should return s2Archive for .tar.sz", func() {
			_, err := FromFilename("filename.tar.sz")
			g.Assert(err == nil).IsTrue("failed to determine .tar.sz suffix")
		})

		g.It("Should return zstdArchive for .tzst", func() {
			_, err := FromFilename("filename.tzst")
			g.Assert(err == nil).IsTrue("failed to determine .tzst suffix")
		})

		g.It("Should return zstdArchive for .tar.zst", func() {
			_, err := FromFilename("filename.tar.zst")
			g.Assert(err == nil).IsTrue("failed to determine .tar.zst suffix")
		})

		g.It("Should return error for everything else", func() {
			_, err := FromFilename("filename.ttt")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("Unknown file format for archive filename.ttt")
		})
	})
}
