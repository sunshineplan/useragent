package useragent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sunshineplan/useragent/internal/verhist"
)

const baseURL = "https://cdn.jsdelivr.net/gh/sunshineplan/useragent/%s"

var cache sync.Map

type value struct {
	useragent string
	expires   time.Time
}

// LatestByOS returns the latest Chrome user agent string for the specified operating system.
func LatestByOS(platform Platform) (string, error) {
	if platform.Normalize(); platform == "" {
		platform = Windows
	}
	if res, ok := cache.Load(platform); ok {
		if v := res.(value); v.expires.Before(time.Now()) {
			cache.Delete(platform)
		} else {
			return v.useragent, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	useragent, err := verhist.UserAgent(ctx, platform.String(), "stable")
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(baseURL, platform.String()), nil)
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
		useragent = string(b)
	}

	cache.Store(platform, value{useragent, time.Now().AddDate(0, 1, 0)})

	return useragent, nil
}

// Latest returns the latest Chrome user agent string for current operating system.
func Latest() (string, error) {
	return LatestByOS(parseGOOS())
}

// UserAgent gets latest chrome user agent string, if failed to get string or
// string is empty, the fallback string will be used.
func UserAgent(fallback string) string {
	ua, err := Latest()
	if err != nil || ua == "" {
		return fallback
	}
	return ua
}
