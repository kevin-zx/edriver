package mywebdriver

import "github.com/tebeka/selenium"

// NewChromeService is using to new a selenium.Service in background
// param: chromeDriverPath is selenium executable binary path, example: chromedriver.exe
// param: serverPort is service servered port
// return: selenium.Service
func NewChromeService(chromeDriverPath string, serverPort int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{
		// 默认selenium的输出都禁止
		selenium.Output(nil),
	}
	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, serverPort, opts...)

	if err != nil {
		return nil, err
	}
	return service, nil

}
