package filex // Package filex decompress used to decompress file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// DecompressFile decompress file
// support ext .zip,.tag.gz
func DecompressFile(archive, dest, filename string) error {
	if exist, _ := PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	if exist, _ := PathExists(dest); !exist {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
	}

	np := fmt.Sprintf("%s/%s", dest, filename)
	if exist, _ := PathExists(np); exist {
		fmt.Println("remove")
		_ = os.RemoveAll(np)
	}

	// guess file type
	// .zip or .tar.gz
	t := filepath.Ext(archive)
	switch t {
	case ".gz":
		err := unTarGz(archive, dest, filename)
		if err != nil {
			return errors.Wrapf(err, "decompress .tar.gz [%s] to dest [%s]", archive, dest)
		}
		return nil
	case ".zip":
		return nil
	default:
		return fmt.Errorf("unsupport file extension [%s]", t)
	}
}

func unTarGz(archive, dest, filename string) error {
	srcFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	hdr, err := tr.Next()
	if err != nil {
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	if !hdr.FileInfo().IsDir() {
		// create path before create file in <create> func, continue here
		return fmt.Errorf("compressed file must contains a dir in the root path")
	}

	rootName := hdr.Name

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		hdrName := hdr.Name
		if !strings.HasPrefix(hdrName, rootName) {
			return fmt.Errorf("compressed file must in a root file while file=%s not in root file=%s", hdrName, rootName)
		}
		hdrName = path.Join(filename, hdrName[len(rootName):])

		filename := path.Join(dest, hdrName) // dest + hdr.Name

		if hdr.FileInfo().IsDir() {
			continue
		}

		file, err := create(filename)
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			_ = file.Close()
			return err
		}
		_ = os.Chmod(filename, hdr.FileInfo().Mode())
		_ = file.Close()
	}

	return nil
}
