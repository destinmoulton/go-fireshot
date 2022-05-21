// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	url := "https://caltopo.com/map.html#ll=36.09336,-105.40867&z=12&b=mbt&a=modis_mp"

	filename := generateScreenshotFilename("caltopo-sipapu")
	path := filepath.Join("screenshots", "caltopo-sipapu", filename)
	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {

		log.Fatal(err)
	}
	if err := ioutil.WriteFile(path, buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	//str := kb.Control + "o"
	return chromedp.Tasks{
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(urlstr),
		chromedp.Sleep(time.Second * 2),
		chromedp.Evaluate("let cookie = document.getElementById('cookie-banner'); cookie.parentNode.removeChild(cookie)", nil),
		chromedp.Evaluate("let node = document.getElementById('page_top'); node.parentNode.removeChild(node)", nil),
		chromedp.Evaluate("let left = document.getElementById('page_left'); left.parentNode.removeChild(left)", nil),
		chromedp.Evaluate("let right = document.getElementById('page_right'); right.parentNode.removeChild(right)", nil),
		chromedp.FullScreenshot(res, quality),
	}
}

func generateScreenshotFilename(prefix string) string {
	dt := time.Now()
	timestamp := dt.Format("D01_02_2006-T15_04_05")
	return fmt.Sprintf("%s-%s.png", prefix, timestamp)
}
