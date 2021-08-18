package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

var api, filename, agent string

func getAgent() error {
	soup.Header("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	body, err := soup.Get(api)
	if err != nil {
		return err
	}

	ua := soup.HTMLParse(body).Find("span", "class", "code")
	if ua.Error != nil {
		return ua.Error
	}

	if agent = ua.Text(); !strings.Contains(agent, "Chrome") {
		return fmt.Errorf("bad result: %s", agent)
	}

	return nil
}

func main() {
	flag.StringVar(&api, "url", "", "URL")
	flag.StringVar(&filename, "filename", "user-agent", "filename")
	flag.Parse()

	if err := getAgent(); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filename, []byte(agent), 0666); err != nil {
		log.Fatal(err)
	}
}
