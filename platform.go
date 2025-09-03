package useragent

import (
	"runtime"
	"slices"
	"strings"
)

const (
	Windows   Platform = "win"
	Windows64 Platform = "win64"
	Mac       Platform = "mac"
	MacARM64  Platform = "mac_arm64"
	Linux     Platform = "linux"
	Android   Platform = "android"
	WebView   Platform = "webview"
	IOS       Platform = "ios"
	Lacros    Platform = "lacros"
)

type Platform string

func (p Platform) String() string { return string(p) }

func SupportedPlatforms() []Platform {
	return []Platform{
		Windows,
		Windows64,
		Mac,
		MacARM64,
		Linux,
		Android,
		WebView,
		IOS,
		Lacros,
	}
}

func (p *Platform) Normalize() {
	*p = Platform(strings.ToLower(string(*p)))
	if !slices.Contains(SupportedPlatforms(), *p) {
		*p = Platform("")
	}
}

func parseGOOS() Platform {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "amd64" {
			return Windows64
		}
		return Windows
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return MacARM64
		}
		return Mac
	case "linux":
		return Linux
	case "android":
		return Android
	case "ios":
		return IOS
	default:
		return ""
	}
}
