package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func DoesExist(f string) bool {
	if stat, err := os.Stat(f); err == nil {
		return !stat.IsDir()
	}
	return false
}

func DoesNotExist(f string) bool  {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return true
	}
	return false
}

func WriteFile(s, f string) error {
	return ioutil.WriteFile(f, []byte(s), os.ModePerm)
}

func ReadFile(f string) (string, error) {
	b, err := ioutil.ReadFile(f)
	return string(b), err
}

func writeStruct(v interface{}, p string) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	dir := path.Dir(p)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err = ioutil.WriteFile(p, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func readStruct(v interface{}, p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}
