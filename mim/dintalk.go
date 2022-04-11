package mim

import (
	"common_go/standard"
	"fmt"
	"github.com/blinkbean/dingtalk"
	"log"
)

// @title    SendTextMsg
// @description    推送文本信息到钉钉群
// @param content 推送的文本内容,token 钉钉群的token,secret 二次加密
// @return

func SendTextMsg(content, token, secret string) error {
	cli := dingtalk.InitDingTalkWithSecret(token, secret)
	err := cli.SendTextMessage(content)
	if err != nil {
		log.Printf("SendTextMsg error : %v\n", err)
		return err
	}
	return nil
}

// @title    SendMarkDownMsg
// @description     推送包含图片链接的makrdown信息到钉钉群
// @param title 标题,text markdown格式的内容(可以包含image url等)
// @return

func SendMarkDownMsg(title, text, token, secret string) error {
	cli := dingtalk.InitDingTalkWithSecret(token, secret)
	err := cli.SendMarkDownMessage(title, text)
	if err != nil {
		standard.XWarning(fmt.Sprintf("SendMarkDownMsg error : %v", err))
		return err
	}
	return nil
}
