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

	html, err := fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := parse(html); err != nil {
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

func fetch() (html node.Node, err error) {
	req, err := http.NewRequest("GET", *api, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	html, err = node.Parse(resp.Body)
	return
}

func parse(html node.Node) error {
	for k, v := range list {
		if h2 := html.Find(0, node.H2, node.String(regexp.MustCompile(v))); h2 != nil {
			if span := h2.Find(node.Next, node.Span, node.Class("code")); span != nil {
				ua := span.GetText()
				if err := os.WriteFile(k, []byte(ua), 0644); err != nil {
					return err
				}
				list[k] = ua
			} else {
				return errors.New("no found")
			}
		} else {
			if h1 := html.Find(0, node.H1); h1 != nil {
				return errors.New(h1.GetText())
			}
			return errors.New("blocked")
		}
	}
	return nil
}
