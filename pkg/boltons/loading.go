package boltons

import (
	"io/ioutil"
	"path/filepath"
)

func LoadLibrary(path string) (out Library, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return out, err
	}

	out = Library{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		out[f.Name()], err = loadBoltons(filepath.Join(path, f.Name()))
		if err != nil {
			break
		}
	}
	return out, err
}

func loadBoltons(path string) ([]Bolton, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	out := []Bolton{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		out = append(out, New(filepath.Join(path, f.Name())))
	}
	return out, nil
}
