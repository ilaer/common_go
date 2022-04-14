package unit_test

import (
	"common_go/network"
	"fmt"
	"gcore/core"
	"testing"
)

func TestHTTPPost(t *testing.T) {
	siteURL := "https://lostcloud.top/auth/login"
	headers := map[string]string{
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36",
		"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	data := map[string]interface{}{
		"email":  " ila2002@qq.com",
		"passwd": " 251405",
		"code":   "",
	}
	httpProxy := "http://127.0.0.1:1080"

	_, err := network.HTTPPost(siteURL, httpProxy, headers, data, 30)
	if err != nil {
		core.XWarning(fmt.Sprintf("HTTPPost error : %v", err))

	}

}
