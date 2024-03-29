package useragent

import (
	"strings"
	"testing"
)

func TestUserAgentString(t *testing.T) {
	ua, err := Latest()
	if err != nil {
		t.Fatal(err)
	}
	if expect := "Chrome"; !strings.Contains(ua, expect) {
		t.Fatalf("expected contains %q; got %q", expect, ua)
	}
}
