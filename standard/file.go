package standard

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// Copy 复制单个文件或整个文件夹
func Copy(from, to string) error {
	f, e := os.Stat(from)
	if e != nil {
		return e
	}
	if f.IsDir() {
		//from是文件夹，那么定义to也是文件夹
		if list, e := ioutil.ReadDir(from); e == nil {
			for _, item := range list {
				if e = Copy(filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else {
		//from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e = os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}
		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建一个文件用于保存
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, e = io.Copy(out, bufReader)
	}
	return e
}

//@Title: ListFilePaths
//@Description: list all filepath from parent path

func ListFilePaths(rootDir, suffix string, PTHS *[]string) error {
	fs, _ := ioutil.ReadDir(rootDir)
	for _, f := range fs {
		if f.IsDir() {
			ListFilePaths(filepath.Join(rootDir, f.Name()), suffix, PTHS)
		} else {
			ext := filepath.Ext(f.Name())
			ext = strings.ToLower(ext)
			if ext != suffix {
				//standard.XWarning(fmt.Sprintf("lifewood#%s", filepath.Join(rootDir, f.Name())))
				continue
			}
			*PTHS = append(*PTHS, filepath.Join(rootDir, f.Name()))
		}
	}
	return nil
}

//@Title: GenerateOutputFilePath
//@Description: create output filepath string from input filepath

func GenerateOutputFilePath(inputPath, outputPath string, FilePaths []string) (map[string]string, error) {
	inputPathName := filepath.Base(inputPath)
	//fmt.Printf("inputPath===>%v\n", inputPath)
	//fmt.Printf("inputPathName===>%v\n", inputPathName)
	outputFilePaths := map[string]string{}
	for _, fp := range FilePaths {
		ret1 := strings.Split(fp, inputPath)
		//fmt.Printf("%v\n", ret1)
		if len(ret1) < 2 {
			continue
		}

		outputFilePath := filepath.Join(outputPath, inputPathName, ret1[1])
		//fmt.Printf("outputFilePath====>%v\n", outputFilePath)
		retPath := filepath.Dir(outputFilePath)
		//fmt.Printf("retPath====>%v\n", retPath)

		exist, _ := FileExists(retPath)
		if exist == false {
			err := os.MkdirAll(retPath, os.ModePerm)
			if err != nil {
				XWarning(fmt.Sprintf("MkdirAll %s error : %v", retPath, err))
				continue
			}
		}

		outputFilePaths[fp] = outputFilePath
	}
	return outputFilePaths, nil
}
