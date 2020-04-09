package main

import (
	"context"
	"go-recruitment-spider/crawl"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	interrupt := make(chan os.Signal, 100)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-interrupt
		cancel()
	}()

	crawl.NewLiepinCrawl(interrupt).Crawl(ctx)
}
