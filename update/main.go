package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/sunshineplan/chrome"
)

func main() {
	ua := chrome.UserAgent()
	fmt.Print(ua)
	if err := os.WriteFile(runtime.GOOS, []byte(ua), 0644); err != nil {
		log.Fatal(err)
	}
}
