package network

import (
	"common_go/im"
	"common_go/standard"
	"fmt"
	"gcore/core"
	"strings"
	"time"
)

// @Title  AccessWebsite
// @Description  定时检测是否可以访问网站,如果第一次正常,发一次正常信息到群里,直到超时时,再发一次信息到群里;如果第一次超时,发一次超时信息到群里,直到恢复正常,再发一次正常信息到群里.
// @param     siteUrl 网站网址,method get或post,httpProxy 代理地址,headers 头部,loop 循环次数,delay 间隔延迟,timeout 超时
// @return   ok 是否可以访问, error        是否成功

func AccessWebsite(siteUrl, method, httpProxy string, headers map[string]string, loop, delay, timeout int, dingTalkToken, dingTalkSecret string) (ok bool, err error) {
	ticker := time.NewTicker(time.Duration(delay) * time.Second)

	normalSend := 0 //正常信息发送计数
	errorSend := 0  //错误消息发送计数

	needSend := false //是否需要发送

	total := 1 //循环计数
	for {
		if total > loop {
			break
		}
		t := <-ticker.C
		standard.XWarning(fmt.Sprintf("current time is :%v", t))
		method = strings.ToLower(method)
		if method == "get" {
			_, err := standard.HttpGet(siteUrl, httpProxy, headers, timeout)
			if err != nil {
				core.XWarning(fmt.Sprintf("%v", err))

				//已发送超时信息,计数大于一,无需发送
				if errorSend > 0 {
					needSend = false
				} else
				//未发送超时信息
				{
					needSend = true
					errorSend += 1
				}

			} else {

				//如果之前是超时,且已经发了一次超时信息;现在正常了,重置超时次数和正常信号.
				if normalSend < 1 {
					needSend = true

					normalSend += 1
					errorSend = 0

				} else {
					needSend = false

				}

			}

			if needSend == true {

				var content string
				if err == nil {
					content = fmt.Sprintf(" http get `%s`  success", siteUrl)
				} else {
					content = fmt.Sprintf(" %v", err)
				}

				//发送信息到钉钉群
				err = im.SendTextMsg(
					content,
					dingTalkToken,
					dingTalkSecret,
				)
				if err != nil {
					core.XWarning(fmt.Sprintf("error : %v", err))
				}
			}

		}

		total++
	}
	return
}
