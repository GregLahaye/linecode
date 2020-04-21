package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func SaveStruct(filename string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	dir := path.Dir(filename)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func LoadStruct(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func QuestionFilename(id int) string {
	return path.Join(questionsDirectory, IntToString(id)+".json")
}

func SplitFilename(filename string) (int, string) {
	parts := strings.Split(filename, ".")
	id, _ := strconv.Atoi(parts[0])
	slug := parts[1]
	return id, slug
}

func CacheDir(filename string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, project, filename), nil
}

func CacheStore(filename string, v interface{}) error {
	filename, err := CacheDir(filename)
	if err != nil {
		return err
	}

	return SaveStruct(filename, v)
}

func CacheRetrieve(filename string, v interface{}) error {
	filename, err := CacheDir(filename)
	if err != nil {
		return err
	}

	return LoadStruct(filename, v)
}

func CacheDestroy(filename string) error {
	filename, err := CacheDir(filename)
	if err != nil {
		return err
	}

	return os.RemoveAll(filename)
}
