package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/sunshineplan/node"
	"github.com/sunshineplan/useragent"
)

var list = map[string]string{
	"windows": "on Windows 10",
	"darwin":  "on macOS",
	"linux":   "on Linux",
	"ios":     "on iOS",
	"android": "on Android",
}

var api = flag.String("url", "", "URL")

func main() {
	flag.Parse()

	if err := update(); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, i := range useragent.SupportedOS() {
		fmt.Fprintf(f, "## %s\n", i)
		fmt.Fprintln(f, "```")
		fmt.Fprintln(f, list[i])
		fmt.Fprintln(f, "```")
	}
}

func update() error {
	req, err := http.NewRequest("GET", *api, nil)
	if err != nil {
		return err
	}
	//req.Header.Set("X-Forwarded-For", "8.8.8.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := node.Parse(resp.Body)
	if err != nil {
		return err
	}
	for k, v := range list {
		if h2 := doc.Find(0, node.Tag("h2"), node.String(regexp.MustCompile(v))); h2 != nil {
			if span := h2.Find(node.Next, node.Span, node.Class("code")); span == nil {
				return errors.New("no found")
			} else {
				ua := span.GetText()
				if err := os.WriteFile(k, []byte(ua), 0644); err != nil {
					return err
				}
				list[k] = ua
			}
		}
	}
	return nil
}
