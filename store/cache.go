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
	return writeStruct(v, p)
}

func RemoveFromCache(f string) error {
	p := path.Join(CacheDir(), f)
	return os.Remove(p)
}

func ReadFromCache(v interface{}, f string) error {
	p := path.Join(CacheDir(), f)
	return readStruct(v, p)
}
