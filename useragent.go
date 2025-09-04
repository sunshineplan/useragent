package useragent

import (
	"context"
	"sync"
	"time"

	"github.com/sunshineplan/useragent/internal/verhist"
)

const baseURL = "https://cdn.jsdelivr.net/gh/sunshineplan/useragent/"

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

		if err := verhist.Grab(ctx, baseURL+platform.String(), &useragent); err != nil {
			return "", err
		}
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
