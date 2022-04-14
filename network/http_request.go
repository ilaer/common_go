// @Author  ila

package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcore/core"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

// @Title  HttpGet
// @Description  http get请求封装
// @param     siteUrl  网址,httpProxy 代理,headers,timeout
// @return   body 网站内容, error        是否成功

func HTTPGet(siteUrl, httpProxy string, headers map[string]string, timeout int) (body []byte, err error) {

	transport := &http.Transport{}

	//判断是否使用代理
	proxy, err := url.Parse(httpProxy)
	if len(httpProxy) > 3 && err == nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	//初始化client
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	//创建get请求
	req, _ := http.NewRequest("GET", siteUrl, nil)

	//请求加入头部
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	//发出get请求
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

func HTTPPost(siteURL, httpProxy string, headers map[string]string, data map[string]interface{}, timeout int) (body []byte, err error) {

	transport := &http.Transport{}

	//判断是否使用代理
	proxy, err := url.Parse(httpProxy)
	if len(httpProxy) > 3 && err == nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	//初始化client
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	contentType := headers["Content-Type"]
	if contentType == "" {
		contentType = headers["content-type"]
	}
	//根据content-type判断数据格式化
	var postData io.Reader
	//判断data是否有数据
	if len(data) > 0 {
		if strings.Contains(contentType, "json") == true {
			rawData, err := json.Marshal(data)
			if err != nil {
				core.XWarning(fmt.Sprintf("json.Marshal error : %v", err))
			}
			postData = bytes.NewReader(rawData)
		} else {
			values := url.Values{}
			for k, v := range data {
				values.Set(k, v.(string))
			}

			postData = strings.NewReader(values.Encode())
		}
	}

	//创建post请求
	req, err := http.NewRequest("POST", siteURL, postData)
	if err != nil {
		core.XWarning(fmt.Sprintf("http.NewRequest error : %v", err))
		return
	}

	//设置post请求headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	//发出post请求
	resp, err := client.Do(req)
	if err != nil {
		core.XWarning(fmt.Sprintf("client.Do error : %v\n", err))
		return
	}

	//读取返回页面数据
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		core.XWarning(fmt.Sprintf("ioutil.ReadAll error : %v", err))
		return
	}
	//把服务器返回的cookie写入本地文件

	cookies := resp.Header.Values("Set-Cookie")
	if len(cookies) > 0 {
		fmt.Printf("%v\n", cookies)
		cookieData, err := json.Marshal(cookies)

		if err != nil {
			core.XWarning(fmt.Sprintf(" json.Marshal error : %v", err))
		} else {
			host := resp.Request.URL.Hostname()
			host = strings.ReplaceAll(host, ".", "_")

			cookieFileName := fmt.Sprintf("%s.json", host)
			err = ioutil.WriteFile(cookieFileName, cookieData, os.ModePerm)
			if err != nil {
				core.XWarning(fmt.Sprintf("ioutil.WriteFile error : %v", err))

			}
		}

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
