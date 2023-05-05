package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func run() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	navigate := chromedp.Navigate(`https://www.google.com/`)

	fellingLuckyCSSSelector := `//div[contains(@class,"FPdoLc lJ9FBc")]//input[contains(@class,"RNmpXc")]`
	clickFellingLucky := chromedp.Click(fellingLuckyCSSSelector, chromedp.NodeReady)

	for _, next := range []chromedp.Action{
		navigate,
		clickFellingLucky,
	} {
		err := chromedp.Run(ctx, next)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("done")
}

func main() {
	run()
}
