package useragent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const url = "https://cdn.jsdelivr.net/gh/sunshineplan/useragent/%s"

func SupportedOS() []string { return []string{"windows", "darwin", "linux", "ios", "android"} }

var cache sync.Map

// LatestByOS returns the latest Chrome user agent string for the specified operating system.
func LatestByOS(os string) (string, error) {
	var supported bool
	for _, i := range SupportedOS() {
		if os == i {
			supported = true
			break
		}
	}
	if !supported {
		os = "windows"
	}
	if res, ok := cache.Load(os); ok {
		return res.(string), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(url, os), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user agent string: %s", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != 200 {
		return "", fmt.Errorf("no StatusOK response: %d", code)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	cache.Store(os, string(b))

	return string(b), nil
}

// Latest returns the latest Chrome user agent string for current operating system.
func Latest() (string, error) {
	return LatestByOS(runtime.GOOS)
}

// UserAgent gets latest chrome user agent string, if failed to get string or
// string is empty, the default string will be used.
func UserAgent(defaultUserAgentString string) string {
	ua, err := Latest()
	if err != nil || ua == "" {
		ua = defaultUserAgentString
	}
	return ua
}
