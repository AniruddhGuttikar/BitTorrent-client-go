package torrentfile

import (
	"os"
	//"github.com/jackpal/bencode-go"
)

func ReadFile(filename string) (string, error) {
	FileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(FileContent), nil

}
