package automate

import (
	"fmt"
	"gcore/core"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
)

type CDriver struct {
	ChromeDriverService *selenium.Service  `json:"chrome_driver_service"`
	WebDriver           selenium.WebDriver `json:"web_driver"`
}

// @title ChromeDriverServiceStart
// @description 开启服务
// @param
// @return
func (m *CDriver) StartService(chromeDriverPath string, chromeDriverPort int) (err error) {
	opts := []selenium.ServiceOption{}

	//selenium.SetDebug(true)

	m.ChromeDriverService, err = selenium.NewChromeDriverService(chromeDriverPath, chromeDriverPort, opts...)
	if nil != err {
		log.Printf("ChromeDriverServiceStart error:%v\n", err)
		core.LogWrite(fmt.Sprintf("ChromeDriverServiceStart error:%v\n", err))
		return
	} else {
		log.Println("ChromeDriverServiceStart success.")
	}

	return
}

// @title ChromeDriverServiceStop
// @description 关闭服务
// @param
// @return
func (m *CDriver) StopService() error {
	//注意这里，server关闭之后，chrome窗口也会关闭
	err := m.ChromeDriverService.Stop()
	if err != nil {
		log.Printf("ChromeDriverServiceStop error:%v\n", err)
		core.LogWrite(fmt.Sprintf("ChromeDriverServiceStop error:%v\n", err))
		return err
	}

	return nil

}

// @title WebDriverStart
// @description 开启webdriver,和新增tab
// @param
// @return
func (m *CDriver) StartDriver(chromeDriverPort int) (WebDriver selenium.WebDriver, err error) {
	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
		//https://github.com/tebeka/selenium/issues/233
		"chromeOptions": map[string]interface{}{
			"excludeSwitches": [1]string{"enable-automation"},
		},
	}
	prefCaps := map[string]interface{}{}
	chromeCaps := chrome.Capabilities{
		Prefs: prefCaps,
		Path:  "",
		Args: []string{
			"--headless",        // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--kiosk",           // 加载启动项页面全屏效果，相当于F11
			"--start-maximized", //  最大化运行（全屏窗口）,不设置，取元素会报错
			"--window-size=1920x1080",
			"--disable-infobars", //  关闭左上方Chrome 正受到自动测试软件的控制的提示
			//"--no-sandbox",// 沙盒,linux下要关闭沙盒模式.
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	// 调起chrome浏览器
	WebDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chromeDriverPort))
	if err != nil {
		//启用--user-data-dir参数后,同时只能打开一个chrome程序.
		//如果已打开chrome程序,必须先关闭已打开的chrome,再运行程序
		fmt.Printf("WebDriverStart connect to the webDriver error : %v\n", err)
		core.LogWrite(fmt.Sprintf("WebDriverStart connect to the webDriver  error : %v", err))
		return
	} else {
		log.Println("connect to the webDriver success.")
	}

	return
}

// @title ChromeDriverServiceStop
// @description 关闭webdriver和打开的窗口
// @param
// @return
func (m *CDriver) StopDriver() error {
	//关闭一个webDriver会对应关闭一个chrome窗口
	//但是不会导致seleniumServer关闭
	err := m.WebDriver.Quit()
	if err != nil {
		log.Printf("WebDriverStop error:%v\n", err)
		core.LogWrite(fmt.Sprintf("WebDriverStop error:%v\n", err))
		return err
	}

	return nil
}
