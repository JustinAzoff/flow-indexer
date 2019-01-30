package backend

import (
	"compress/bzip2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	gzip "github.com/klauspost/pgzip"
)

//PipedDecompressor is used to wrap higher performing native decompression tools
type PipedDecompressor struct {
	io.ReadCloser
	wrapped io.ReadCloser
	cmd     *exec.Cmd
}

func NewPipedDecompressor(r io.ReadCloser, prog string) (*PipedDecompressor, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = r
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	pd := &PipedDecompressor{
		ReadCloser: out,
		wrapped:    r,
		cmd:        cmd,
	}
	return pd, err
}

func (pd *PipedDecompressor) Close() error {
	pd.cmd.Wait()
	pd.wrapped.Close()
	pd.ReadCloser.Close()
	return nil
}

type WrappedDecompressor struct {
	io.ReadCloser
	wrapped io.ReadCloser
}

func (wd *WrappedDecompressor) Close() error {
	wd.wrapped.Close()
	return wd.ReadCloser.Close()
}

func OpenDecompress(fn string) (r io.ReadCloser, err error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(fn)

	switch ext {
	case ".log", ".txt":
		return f, err
	case ".gz":
		gzr, err := gzip.NewReader(f)
		return &WrappedDecompressor{
			ReadCloser: gzr,
			wrapped:    f,
		}, err
	case ".bz2":
		bzr := bzip2.NewReader(f)
		return &WrappedDecompressor{
			ReadCloser: ioutil.NopCloser(bzr),
			wrapped:    f,
		}, nil
	case ".xz":
		xzr, err := NewPipedDecompressor(f, "xzcat")
		return xzr, err
	default:
		return f, err
	}
}
