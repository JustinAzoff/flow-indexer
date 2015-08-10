package main

import (
	//"compress/bzip2"
	//"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type PipedDecompressor struct {
	io.ReadCloser
	wrapped io.ReadCloser
	cmd     *exec.Cmd
}

func NewPipedDecompressor(f *os.File, prog string) (*PipedDecompressor, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = f
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	pd := &PipedDecompressor{
		ReadCloser: out,
		wrapped:    f,
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
		//return gzip.NewReader(f)
		return NewPipedDecompressor(f, "zcat")
	case ".bz2":
		return NewPipedDecompressor(f, "bzcat")
	}

	return nil, fmt.Errorf("Unknown filetype %s", ext)
}
