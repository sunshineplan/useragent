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

	agent = ua.Text()

	if !strings.Contains(agent, "Chrome") {
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

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(agent)
}
