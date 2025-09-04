package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sunshineplan/useragent"
	"github.com/sunshineplan/useragent/internal/verhist"
)

func main() {
	for _, platform := range useragent.SupportedPlatforms() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		useragent, err := verhist.UserAgent(ctx, platform.String(), "stable")
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(platform, useragent)
		if err := os.WriteFile(platform.String(), []byte(useragent), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
