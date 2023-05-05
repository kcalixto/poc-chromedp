package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

func run() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	done := make(chan string, 1)

	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			completed := "(unknown)"
			if ev.TotalBytes != 0 {
				completed = fmt.Sprintf("%0.2f%%", ev.ReceivedBytes/ev.TotalBytes*100.0)
			}

			log.Printf("state: %s, completed: %s\n", ev.State.String(), completed)

			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}
	})

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	navigate := chromedp.Navigate(`https://github.com/chromedp/examples`)

	clickRepoSummary := chromedp.Click(`//get-repo//summary`, chromedp.NodeReady)

	setDownloadBehavior := browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
		WithDownloadPath(wd).
		WithEventsEnabled(true)

	clickDownloadZip := chromedp.Click(`//get-repo//a[contains(@data-ga-click, "download zip")]`, chromedp.NodeVisible)

	tasks := []chromedp.Action{
		navigate,
		clickRepoSummary,
		setDownloadBehavior,
		clickDownloadZip,
	}

	for _, next := range tasks {
		err = chromedp.Run(ctx, next)
		// Note: Ignoring the net::ERR_ABORTED page error is essential here
		// since downloads will cause this error to be emitted, although the
		// download will still succeed.
		if err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
			log.Fatal(err)
		}
	}

	guid := <-done

	log.Printf("wrote %s", filepath.Join(wd, guid))
}

func main() {
	run()
}
