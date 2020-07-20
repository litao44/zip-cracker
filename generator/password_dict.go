package generator

import (
	"bufio"
	"io"
	"sync"
)

func NewDictPasswordGenerator(r io.Reader) (*DictPasswordGenerator, error) {
	return &DictPasswordGenerator{
		reader:  r,
		scanner: bufio.NewScanner(r),
		lock:    sync.Mutex{},
	}, nil
}

type DictPasswordGenerator struct {
	reader  io.Reader
	scanner *bufio.Scanner
	lock    sync.Mutex
}

func (dp *DictPasswordGenerator) Generate() (string, bool) {
	dp.lock.Lock()
	defer dp.lock.Unlock()

	if dp.scanner.Scan() {
		return dp.scanner.Text(), false
	}
	return "", true
}
