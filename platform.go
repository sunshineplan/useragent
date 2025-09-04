package useragent

import (
	"runtime"
	"slices"
	"strings"
)

// https://versionhistory.googleapis.com/v1/chrome/platforms/
const (
	Windows   Platform = "win"
	Windows64 Platform = "win64"
	Mac       Platform = "mac"
	MacARM64  Platform = "mac_arm64"
	Linux     Platform = "linux"
	Android   Platform = "android"
)

type Platform string

func (p Platform) String() string {
	return strings.ToLower(string(p))
}

func SupportedPlatforms() []Platform {
	return []Platform{
		Windows,
		Windows64,
		Mac,
		MacARM64,
		Linux,
		Android,
	}
}

func (p *Platform) Normalize() {
	*p = Platform(p.String())
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
	//case "ios":
	//	return IOS
	default:
		return ""
	}
}
