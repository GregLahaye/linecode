package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func WriteStruct(v interface{}, p string) error {
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

func ReadStruct(v interface{}, p string) error {
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
