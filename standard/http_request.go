// @Author  ila

package standard

import (
	"fmt"
	"gcore/core"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// @Title  HttpGet
// @Description  http get请求封装
// @param     siteUrl  网址,httpProxy 代理,headers,timeout
// @return   body 网站内容, error        是否成功

func HttpGet(siteUrl, httpProxy string, headers map[string]string, timeOut int) (body []byte, err error) {

	transport := &http.Transport{}

	//判断是否使用代理
	proxy, err := url.Parse(httpProxy)
	if err == nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	//初始化client
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeOut) * time.Second,
	}

	//method加入client
	req, _ := http.NewRequest("GET", siteUrl, nil)

	//请求加入头部
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	//发出请求
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("client.Do error : %v\n", err)
		return
	}

	//读取返回页面数据
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		core.XWarning(fmt.Sprintf("ioutil.ReadAll error : %v", err))
		return
	}
	defer resp.Body.Close()
	return
}

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
