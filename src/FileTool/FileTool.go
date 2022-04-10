package filetool

import (
	"io/ioutil"
	"os"
	"path"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Clear(basePath string) error {
	if ok, _ := PathExists(basePath); !ok {
		return os.Mkdir(basePath, 0777)
	}
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
func WriteFileB(path string, buf []byte) error {
	return ioutil.WriteFile(path, buf, 0777)
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func Home() string {
	return "/code/TMManager/"
	//return os.Getenv("CUGOJ_HOME")
}
