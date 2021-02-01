package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

func getAgent(url string) (string, error) {
	soup.Header("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	body, err := soup.Get(url)
	if err != nil {
		return "", err
	}
	ua := soup.HTMLParse(body).Find("span", "class", "code")
	if ua.Error != nil {
		return "", ua.Error
	}
	if !strings.Contains(ua.Text(), "Chrome") {
		return "", errors.New("Not Chrome")
	}
	return ua.Text(), nil
}

func main() {
	url := flag.String("url", "", "URL")
	flag.Parse()

	agent, err := getAgent(*url)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("chrome.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(agent)
}
