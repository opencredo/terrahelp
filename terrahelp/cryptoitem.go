package terrahelp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type cryptoItemAction func(*CryptoHandlerOpts, CryptoItem) error

// CryptoItem defines actions to validate, read and write
// content required as input to, and output from, a cryptographic action.
type CryptoItem interface {
	validateCryptoItem() error
	readFromSource() ([]byte, error)
	beforeCryptoAction() error
	writeToTarget([]byte) error
}

// StreamCryptoItem defines actions to validate, read and write
// content from a stream (typically stdin and stdout) required as input to,
// and output from, a cryptographic action.
// Note: There generally be NO arbitrary writing to stdout for logging purposes
// with a StreamCryptoItem as often the stream is stdin / stdout itself.
type StreamCryptoItem struct {
	in  io.Reader
	out io.Writer
}

// NewStreamCryptoItem creates a new StreamCryptoItem
func NewStreamCryptoItem(in io.Reader, out io.Writer) *StreamCryptoItem {
	return &StreamCryptoItem{in, out}
}

// NewStdStreamCryptoItem creates a new StreamCryptoItem using
// stdin for obtaining the input into, and stdout for writing
// the result of a cryptographic action
func NewStdStreamCryptoItem() *StreamCryptoItem {
	return NewStreamCryptoItem(os.Stdin, os.Stdout)
}

func (p *StreamCryptoItem) validateCryptoItem() error {
	return nil
}

func (p *StreamCryptoItem) readFromSource() ([]byte, error) {
	return ioutil.ReadAll(p.in)
}

func (p *StreamCryptoItem) writeToTarget(b []byte) error {
	_, err := p.out.Write(b)
	return err
}
func (p *StreamCryptoItem) beforeCryptoAction() error {
	return nil
}

// FileCryptoItem defines actions to validate, read and write
// content from a file on the filesystem required as input to,
// and output from, a cryptographic action.
type FileCryptoItem struct {
	filename string
	bkp      bool
	bkpExt   string
}

// NewFileCryptoItem creates a new FileCryptoItem
func NewFileCryptoItem(f string, bkp bool, bkpExt string) *FileCryptoItem {
	return &FileCryptoItem{f, bkp, bkpExt}
}

func (f *FileCryptoItem) validateCryptoItem() error {
	fi, err := os.Stat(f.filename)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%s must be a valid file", f.filename)
	}
	return nil
}

func (f *FileCryptoItem) readFromSource() ([]byte, error) {
	if err := f.validateCryptoItem(); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(f.filename)
}

func (f *FileCryptoItem) writeToTarget(b []byte) error {
	if err := f.validateCryptoItem(); err != nil {
		return err
	}
	return ioutil.WriteFile(f.filename, b, 0777)
}
func (f *FileCryptoItem) beforeCryptoAction() error {
	if f.bkp {
		bkp := f.filename + f.bkpExt
		log.Printf("Backuping up %s --> %s ", f.filename, bkp)
		return CopyFile(f.filename, bkp)
	}
	return nil
}
