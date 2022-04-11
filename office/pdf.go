package office

import (
	"bytes"
	"fmt"
	"gcore/core"
	"github.com/axgle/mahonia" //编码转换
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// @title ConvertPdfToPng
// @description 把整个pdf渲染成png
// @param pdf filepath
// @return error

func ConvertPdfToPng(pdfPath string, subPath, fileName string) error {
	currentDir, _ := os.Getwd()
	xpdfPath := filepath.Join(currentDir, "tool", "pdftopng.exe")
	//println(xpdfPath)
	_, err := os.Stat(xpdfPath)
	if os.IsNotExist(err) {
		fmt.Println("file no exist.")
		return fmt.Errorf("pdftopng.exe not found in ./tools")
	}

	//cmdLine := fmt.Sprintf("cd %s && %s  %s %s", subPath,xpdfPath,pdfPath, fileName)
	//println(cmdLine)
	cmdArgs := []string{
		"/C",
		"cd",
		subPath,
		"&&",
		xpdfPath,
		"-f",
		"1",
		"-l",
		"1",
		"-r",
		"250",
		pdfPath,
		fileName,
	}
	println(strings.Join(cmdArgs, " "))
	c := exec.Command("cmd.exe", cmdArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var enc mahonia.Decoder
	enc = mahonia.NewDecoder("gbk")
	c.Stdout = &stdout
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		fmt.Printf("XPdfToPng Error: %s %s\n", err, enc.ConvertString(stderr.String()))
		//go Xlog(logPath, strings.Join(cmdArgs, " "))
		go core.LogWrite(fmt.Sprintf("%s,XPdfToPng Error: %s %s\n", pdfPath, err, enc.ConvertString(stderr.String())))
		return fmt.Errorf("%s,XPdfToPng Error: %s %s\n", pdfPath, err, enc.ConvertString(stderr.String()))
	}
	//fmt.Printf("XPdfToPng Result: %s\n", stdout.String())
	return nil
}
