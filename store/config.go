package store

import (
	"github.com/GregLahaye/linecode/linecode"
	"os"
	"path"
)

func ConfigDir() string {
	cd, _ := os.UserConfigDir()
	return path.Join(cd, linecode.Project)
}

func SaveToConfig(v interface{}, f string) error {
	p := path.Join(ConfigDir(), f)
	return WriteStruct(v, p)
}

func ReadFromConfig(v interface{}, f string) error {
	p := path.Join(ConfigDir(), f)
	return ReadStruct(v, p)
}

