package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

var (
	api      = flag.String("url", "", "URL")
	filename = flag.String("filename", "user-agent", "filename")
)

func main() {
	flag.Parse()

	agent, err := getAgent()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*filename, []byte(agent), 0666); err != nil {
		log.Fatal(err)
	}
}

func getAgent() (string, error) {
	soup.Header("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	body, err := soup.Get(*api)
	if err != nil {
		return "", err
	}

	ua := soup.HTMLParse(body).Find("span", "class", "code")
	if ua.Error != nil {
		return "", ua.Error
	}

	agent := ua.Text()
	if !strings.Contains(agent, "Chrome") {
		return "", fmt.Errorf("bad result: %s", agent)
	}

	return agent, nil
}
