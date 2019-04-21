package simplex

import (
	"io/ioutil"
	"os"
)

// LoadLtxFile loads a .ltx file and return its content as string
func LoadLtxFile(filename string) (string, error) {
	file, err := os.Open(filename + ".ltx")
	if err != nil {
		return "", err
	}
	defer file.Close()

	contentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}
