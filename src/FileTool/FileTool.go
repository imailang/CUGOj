package filetool

import (
	"io/ioutil"
	"os"
	"path"
)

func Clear(basePath string) error {
	err := os.RemoveAll(basePath)
	if err != nil {
		return err
	}
	err = os.Mkdir(basePath, os.FileMode(0777))
	return err
}

func ReadCases(basePath string) ([]string, error) {
	fs, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []string{}, err
	}
	mp := make(map[string]int)
	for _, f := range fs {
		if !f.IsDir() {
			if ext := path.Ext(f.Name()); ext == "in" || ext == "out" {
				mp[path.Base(f.Name())]++
			}
		}
	}
	cases := []string{}
	for key, val := range mp {
		if val >= 2 {
			cases = append(cases, key)
		}
	}
	return cases, nil
}

func WriteFile(path string, str *string) error {
	return ioutil.WriteFile(path, []byte(*str), 0777)
}
