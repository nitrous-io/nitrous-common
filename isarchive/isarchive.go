package isarchive

import (
	"archive/tar"
	"bytes"
)

// extracted from https://github.com/docker/docker/blob/master/archive/archive.go

type Compression int

const (
	Uncompressed Compression = iota
	Bzip2
	Gzip
	Xz
)

const MaxHeaderLength = 512

func IsArchive(header []byte) bool {
	compression := DetectCompression(header)
	if compression != Uncompressed {
		return true
	}
	r := tar.NewReader(bytes.NewBuffer(header))
	_, err := r.Next()
	return err == nil
}

func DetectCompression(source []byte) Compression {
	for compression, m := range map[Compression][]byte{
		Bzip2: {0x42, 0x5A, 0x68},
		Gzip:  {0x1F, 0x8B, 0x08},
		Xz:    {0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00},
	} {
		if len(source) < len(m) {
			continue
		}
		if bytes.Compare(m, source[:len(m)]) == 0 {
			return compression
		}
	}
	return Uncompressed
}
