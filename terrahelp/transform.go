package terrahelp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// TransformOpts holds the specific options detailing how, and on what
// the transformation action should be performed.
type TransformOpts struct {
	TransformItems []Transformable
	TfvarsFilename string
}

// Transformable defines the set of actions which can be performed on some underlying
// content as part of applying a transformation process to it.
type Transformable interface {
	validate() error
	read() ([]byte, error)
	beforeTransform() error
	write([]byte) error
}

// StreamTransformable defines the set of actions which can be performed on an underlying
// stream of data (typically stdin and stdout) as part of applying a transformation process to it.
// Note: There generally be NO arbitrary writing to stdout for logging purposes
// with a StreamTransformable as often the stream is stdin / stdout itself.
type StreamTransformable struct {
	in  io.Reader
	out io.Writer
}

// NewStreamTransformable creates a new StreamTransformable
func NewStreamTransformable(in io.Reader, out io.Writer) *StreamTransformable {
	return &StreamTransformable{in, out}
}

// NewStdStreamTransformable creates a new StreamTransformable using
// stdin for obtaining the input into, and stdout for writing
// the result of a transformation action
func NewStdStreamTransformable() *StreamTransformable {
	return NewStreamTransformable(os.Stdin, os.Stdout)
}

func (p *StreamTransformable) validate() error {
	return nil
}

func (p *StreamTransformable) read() ([]byte, error) {
	return ioutil.ReadAll(p.in)
}

func (p *StreamTransformable) write(b []byte) error {
	_, err := p.out.Write(b)
	return err
}
func (p *StreamTransformable) beforeTransform() error {
	return nil
}

// FileTransformable defines actions to validate, read and write
// content from a file on the filesystem required as input to,
// and output from, a transformation action.
type FileTransformable struct {
	filename string
	bkp      bool
	bkpExt   string
}

// NewFileTransformable creates a new FileTransformable
func NewFileTransformable(f string, bkp bool, bkpExt string) *FileTransformable {
	return &FileTransformable{f, bkp, bkpExt}
}

func (f *FileTransformable) validate() error {
	fi, err := os.Stat(f.filename)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%s must be a valid file", f.filename)
	}
	return nil
}

func (f *FileTransformable) read() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(f.filename)
}

func (f *FileTransformable) write(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	return ioutil.WriteFile(f.filename, b, 0777)
}
func (f *FileTransformable) beforeTransform() error {
	if f.bkp {
		bkp := f.filename + f.bkpExt
		log.Printf("Backuping up %s --> %s ", f.filename, bkp)
		return CopyFile(f.filename, bkp)
	}
	return nil
}
