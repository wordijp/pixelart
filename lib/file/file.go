package file

import (
	"bytes"
	"os"
)

// ReadAll -- ファイルを全読み込む
func ReadAll(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
