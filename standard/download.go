package standard

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
)

// Wget 调用wget.exe下载,-t(retries)参数表示重试次数,-T(timeout)参数表示超时等待时间,-w(waitRetry)两次尝试之间间隔SECONDS秒,-O把文档写到FILE文件中
func WGet(url, localPath string, retries, timeout, waitRetry int) (err error) {
	cmd := exec.Command(
		"wget.exe",
		url,
		"-t",
		fmt.Sprintf("%d", retries),
		"-T",
		fmt.Sprintf("%d", timeout),
		"-w",
		fmt.Sprintf("%d", waitRetry),
		"-O",
		localPath,
		"--no-check-certificate", //不检查sll
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {

		fmt.Println(fmt.Sprint(err) + ": " + ConvertByte2String([]byte(stderr.String()), GB18030))
		return err
	}
	fmt.Println("Result: " + stdout.String())

	return
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// ConvertByte2String cmd输出内容转中文
func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
