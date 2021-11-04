package core

// @Title  archive
// @Description  创建和释放压缩包
// @Author  ila 2021-11-04 14:33:00
// @Update  ila

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// @title    Unzip
// @description    解压zip文件到指定路径,只解压文件名包含contains字符串的文件.
// @auth      ila 2021-11-04 14:33:00
// @param     zipPath  zip文件路径,attachPath 解压路径,contains 包含的字符串
// @return    error        是否成功
func Unzip(zipPath, attachPath string, contains []string) error {

	z, err := zip.OpenReader(zipPath)
	defer z.Close()

	if err != nil {
		return err
	}

	for _, f := range z.File {
		if f.FileInfo().IsDir() {
			continue
		}

		//判断文件名包含所需字符
		containFlag := false
		for _, c := range contains {
			if strings.Contains(f.Name, c) == true {
				containFlag = true
				break
			}
		}

		//文件名没有包含所需字符,跳过
		if containFlag == false {
			continue
		}
		path := filepath.Join(attachPath, f.Name)

		fo, err := f.Open()
		defer fo.Close()
		if err != nil {
			log.Printf("zip file Open error :%v\n", err)
			continue
		}
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode())
		defer fw.Close()
		if err != nil {
			continue
		}
		_, err = io.Copy(fw, fo)
		if err != nil {
			log.Printf("zip file write to attach error :%v\n", err)
		}
	}

	return nil
}
