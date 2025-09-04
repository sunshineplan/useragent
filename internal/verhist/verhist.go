// https://github.com/chromedp/verhist
package verhist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultVersion = "140.0.0.0"
	uaFormat       = "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s%s Safari/537.36"
)

var BaseURL = "https://versionhistory.googleapis.com/v1/chrome/platforms/"

// Versions returns the versions for the platform, channel.
func Versions(ctx context.Context, platform, channel string, q ...string) ([]Version, error) {
	if len(q) == 0 {
		q = []string{
			"order_by", "version desc",
		}
	}
	res := new(VersionsResponse)
	if err := Grab(ctx, BaseURL+platform+"/channels/"+channel+"/versions", res, q...); err != nil {
		return nil, err
	}
	return res.Versions, nil
}

// Latest returns the latest version for the platform, channel.
func Latest(ctx context.Context, platform, channel string) (Version, error) {
	versions, err := Versions(ctx, platform, channel)
	switch {
	case err != nil:
		return Version{}, err
	case len(versions) == 0:
		return Version{}, fmt.Errorf("no versions available")
	}
	return versions[0], nil
}

// UserAgent builds the user agent for the platform, channel.
func UserAgent(ctx context.Context, platform, channel string) (string, error) {
	latest, err := Latest(ctx, platform, channel)
	if err != nil {
		return "", err
	}
	return latest.UserAgent(platform), nil
}

// VersionsResponse wraps the versions API response.
type VersionsResponse struct {
	Versions      []Version `json:"versions,omitempty"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

// Version contains information about a chrome release.
type Version struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// UserAgent builds a user agent for the platform.
func (ver Version) UserAgent(platform string) string {
	typ, extra := "Windows NT 10.0; Win64; x64", ""
	switch platform {
	case "linux":
		typ = "X11; Linux x86_64"
	case "mac", "mac_arm64":
		typ = "Macintosh; Intel Mac OS X 10_15_7"
	case "android":
		typ, extra = "Linux; Android 10; K", " Mobile"
	}
	v := defaultVersion
	if i := strings.Index(ver.Version, "."); i != -1 {
		v = ver.Version[:i] + ".0.0.0"
	}
	return fmt.Sprintf(uaFormat, typ, v, extra)
}

// Grab grabs the url and json decodes it.
func Grab(ctx context.Context, urlstr string, v any, q ...string) error {
	if len(q)%2 != 0 {
		return fmt.Errorf("invalid query")
	}
	z := make(url.Values)
	for i := 0; i < len(q); i += 2 {
		z.Add(q[i], q[i+1])
	}
	s := z.Encode()
	if s != "" {
		s = "?" + s
	}
	req, err := http.NewRequestWithContext(ctx, "GET", urlstr+s, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not retrieve %s (status: %d)", urlstr, res.StatusCode)
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
