package resource

import (
	"embed"
	"io"
	"path/filepath"
)

//go:embed root
var FS embed.FS

func ReadPass(name string) ([]byte, error) {
	keyFD, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer keyFD.Close()

	keyData, err := io.ReadAll(keyFD)
	if err != nil {
		return nil, err
	}

	passLen := len(keyData) / 33
	realPass := []byte{}
	for i := 0; i < passLen; i++ {
		realPass = append(realPass, keyData[i*33])
	}
	return realPass, nil
}

func ReadAll(name string) ([]byte, error) {
	keyFD, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer keyFD.Close()

	return io.ReadAll(keyFD)
}

func ReadFilenames(dir string) ([]string, error) {
	dirEntries, err := FS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filenames := []string{}

	for _, entry := range dirEntries {

		if !entry.Type().IsRegular() {
			continue
		}

		filename := filepath.Join(dir, entry.Name())

		filenames = append(filenames, filename)
	}
	return filenames, nil
}
