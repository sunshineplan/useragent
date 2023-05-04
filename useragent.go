package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sunshineplan/node"
)

var (
	api      = flag.String("url", "", "URL")
	filename = flag.String("filename", "README.md", "filename")
)

func main() {
	flag.Parse()

	agent, err := getAgent()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*filename, []byte(agent), 0644); err != nil {
		log.Fatal(err)
	}
}

func getAgent() (string, error) {
	req, err := http.NewRequest("GET", *api, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := node.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	ua := doc.Find(0, node.Span, node.Class("code")).String()
	if ua == nil {
		return "", errors.New("no found")
	}

	agent := ua.String()
	if !strings.Contains(agent, "Chrome") {
		return "", fmt.Errorf("bad result: %s", agent)
	}

	return agent, nil
}
