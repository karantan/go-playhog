package main

import (
	"fmt"
	"log"

	"playhog/browser"

	"github.com/playwright-community/playwright-go"
)

func main() {
	b, err := browser.GetBrowser(false)
	if err != nil {
		fmt.Printf("count not open the browser: %v", err)
		return
	}
	defer b.Close()

	page, err := b.NewPage()
	if err != nil {
		fmt.Printf("could not create new page: %v", err)
		return
	}

	domain := "karantan.si"
	// Navigate to the URL
	_, err = page.Goto("https://" + domain)
	if err != nil {
		log.Fatalf("could not go to the URL: %v", err)
	}
	//
	// Step 1: Load the page, inject posthog and do some action (e.g. click on a link)
	//
	if err := browser.InjectPostHog(page, domain); err != nil {
		log.Fatalf("could not inject PostHog: %v", err)
	}
	// page.Pause()

	locLink := page.GetByRole("link", playwright.PageGetByRoleOptions{Name: "Tips on reading the systemd"})
	locLink.Hover()
	locLink.Click()

	//
	// Step 2: Inject the posthog script and do some action (e.g. fill a form)
	//
	if err := browser.InjectPostHog(page, domain); err != nil {
		log.Fatalf("could not inject PostHog: %v", err)
	}
	locHeader := page.Locator("header").Filter(playwright.LocatorFilterOptions{HasText: "December 3, 2020December 3,"}).GetByRole("link")
	locHeader.Hover()
	locHeader.Click()

	locComment := page.GetByLabel("Comment *")
	locComment.Hover()
	locComment.Click()
	locComment.Fill("test 8")

	// Go back to the previous step (we are on the same page so no injecting required)
	page.GoBack()
	locComment.Hover()
}
