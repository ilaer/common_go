package standard

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// GetChromeVersion 获取chrome的版本

func GetChromeVersion() int {
	ChromeVersion := 90

	//从注册表
	key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Google\\Chrome\\BLBeacon", registry.READ)

	if err != nil {
		XWarning(fmt.Sprintf("GetChromeVersionFromRegedit OpenKey  error : %v\n", err))
		return ChromeVersion
	}
	defer key.Close()

	version, _, err := key.GetIntegerValue("version")
	if err != nil {
		XWarning(fmt.Sprintf("GetChromeVersionFromRegedit OpenKey  error : %v\n", err))
		return ChromeVersion
	}

	ChromeVersion = int(version)

	return ChromeVersion
}
