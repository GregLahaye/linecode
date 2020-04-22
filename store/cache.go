package store

import (
	"github.com/GregLahaye/linecode/linecode"
	"os"
	"path"
)

func CacheDir() string {
	cd, _ := os.UserCacheDir()
	return path.Join(cd, linecode.Project)
}

func SaveToCache(v interface{}, f string) error {
	p := path.Join(CacheDir(), f)
	return WriteStruct(v, p)
}

func ReadFromCache(v interface{}, f string) error {
	p := path.Join(CacheDir(), f)
	return ReadStruct(v, p)
}
