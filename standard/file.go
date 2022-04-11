package standard

import (
	"os"
	"path/filepath"
)

// FileExists 检测文件或目录是否存在
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileNameFromFilePath(path string) (fileName string, err error) {
	exists, err := FileExists(path)
	if exists == false {
		XWarning("file is not exist")
		return
	}
	_, fileNameWithSuffix := filepath.Split(path)
	fileName = filepath.Base(fileNameWithSuffix)

	return
}
