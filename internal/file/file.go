package file

import (
	"os"
)

func CreateFile(filePath string) error {
	_, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteFile(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
