package standard

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// GetContentWithCookies 带cookie获取源地址的内容
func GetContentWithCookies(srcUrl string, cookies []*http.Cookie) (string, error) {

	var content string

	//new个cookiejar罐子
	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Printf("cookiejar New error : %v\n", err)
		return content, err
	}

	//原生的url
	url, _ := url.Parse(srcUrl)

	//把传入的cookies装入罐
	jar.SetCookies(url, cookies)

	//初始化请求
	client := &http.Client{Jar: jar, Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", srcUrl, nil)

	//请求加入头部
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")

	//发出请求
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("client.Do error : %v\n", err)
		return content, err
	}

	//读取返回页面数据
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	content = string(body)

	return content, nil
}
