package mywebdriver

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"
	"time"
)

type EChromeDriver struct {
	selenium.WebDriver
	userAgent     string
	disableWindow bool
}
type DeviceName string

// MobileDeviceName must period update
// sometime chrome update those device will not be supported
const (
	DeviceGalaxyS5 = "Galaxy S5"
	DevicePixel2   = "Pixel 2"
	DeviceIPhoneX  = "iPhone X"
)

// NewChromeWebDriver new a chrome driver with a exist selenium chrome service
// servicePort: a chrome service's port, user edriver service to start
// proxy: start browser with proxy, proxy like protocol://host: port, if the proxy is empty, will not use the proxy
// mobileDeviceName: if you want to start a chrome driver in a mobile device mode you can set,
//		I provided server deviceName.but it's maybe removed by chrome defaults device, the more general implement is specified resolution and user-agent, or get chrome default device list
// UserAgent: specify User-agent, if User-agent is empty use chrome default User-agent
// enableLog: enableLog can enable performance log you can use it to intercept requests, but can't intercept response (or I can't), default false
// disableWindow: if true will set --headless argument to chrome Capabilities chrome will running in the background
// ps. this function has too many params, maybe can implement by options mode, maybe I'll do it
func NewChromeWebDriver(servicePort int, proxy string, mobileDeviceName DeviceName, UserAgent string, enableLog bool, disableWindow bool) (EChromeDriver, error) {
	wdu := EChromeDriver{}
	err := wdu.createWebDriver(servicePort, proxy, mobileDeviceName, enableLog)
	if err != nil {
		wdu.QuitDriver()
	}
	wdu.userAgent = UserAgent
	wdu.disableWindow = disableWindow
	return wdu, err
}

func (myWd *EChromeDriver) createWebDriver(servicePort int, proxy string, mobileDeviceName DeviceName, enableLog bool) (err error) {

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	if enableLog {
		caps.SetLogLevel(log.Performance, log.Info)
	}

	c := chrome.Capabilities{
		ExcludeSwitches: []string{"enable-automation"},
		Args: []string{
			//"--headless",
			"--no-sandbox",
			"--disable-gpu-sandbox",
		},
		W3C: false,
	}
	if myWd.disableWindow {
		c.Args = append(c.Args, "--headless")
	}
	if myWd.userAgent != "" {
		EnableNetworkFlg := true
		c.Args = append(c.Args, "--user-agent="+myWd.userAgent)
		c.PerfLoggingPrefs = &chrome.PerfLoggingPreferences{
			EnableNetwork: &EnableNetworkFlg,
		}
	}
	if mobileDeviceName != "" {
		c.MobileEmulation = &chrome.MobileEmulation{
			DeviceName: string(mobileDeviceName),
		}
	}
	if proxy != "" {
		c.Args = append(c.Args, fmt.Sprintf("--proxy-server=%s", proxy))
	}
	caps.AddChrome(c)
	myWd.WebDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", servicePort))
	if err != nil {
		return err
	}
	err = myWd.WebDriver.SetPageLoadTimeout(20 * time.Second)
	return
}

func (myWd *EChromeDriver) QuitDriver() {
	if myWd.WebDriver != nil {
		_ = (myWd.WebDriver).Close()
		_ = myWd.WebDriver.Quit()
	}
}

// CloseOther Windows sometimes we execute many actions we don't know where we are, and how to go back the method can reset chrome tabs, keep only one tab, and switch to it.
func (myWd *EChromeDriver) CloseOtherWindows() error {
	cw, err := myWd.CurrentWindowHandle()
	if err != nil {
		return err
	}
	ws, err := myWd.WindowHandles()
	if err != nil {
		return err
	}
	for _, w := range ws {
		if w != cw {
			err = myWd.CloseWindow(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
