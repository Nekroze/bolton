package boltons

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Stores the path to the bolton file
type Bolton struct {
	Path string
	Name string
}

func (b Bolton) String() string {
	return b.Name
}

func (b Bolton) Contents() (string, error) {
	bs, err := ioutil.ReadFile(b.Path)
	return string(bs), err
}

func New(path string) Bolton {
	return Bolton{
		Path: path,
		Name: filepath.Base(path),
	}
}

// Array of boltons associated with their tag for filtering
type Library map[string][]Bolton

func (bl Library) Get(address string) (out Bolton, err error) {
	parts := strings.Split(address, ":")
	for _, bolton := range bl[parts[0]] {
		if bolton.String() == parts[1] {
			return bolton, nil
		}
	}
	return out, fmt.Errorf("no bolton found called %s under the tag %s", parts[0], parts[1])
}
