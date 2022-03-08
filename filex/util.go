package filex

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rdsutbbp/utilx"
	"os"
	"path/filepath"
)

func Create(name string) (*os.File, error) {
	dir, _ := filepath.Split(name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

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

// Open file
func Open(name string) (file *os.File, err error) {
	splits, _ := filepath.Split(name)
	if err := os.MkdirAll(splits, 0755); err != nil {
		return nil, err
	}

	output, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// Ext return file extension
// e.g. Ext("abc.json") => "json"
func Ext(filename string, checkList ...string) (ext string, err error) {
	ext = filepath.Ext(filename)
	if len(ext) <= 1 {
		return ext, fmt.Errorf("filename: %s requires valid extension", filename)
	}

	ext = ext[1:]
	if len(checkList) > 0 && !utilx.StringInSlice(ext, checkList) {
		return ext, errors.Errorf("Unsupported Config Type %s", ext)
	}

	return ext, nil
}