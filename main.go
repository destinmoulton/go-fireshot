// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	config := loadJSONConfig("config/config.json")
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	for _, shot := range config.Shots {

		filename := generateScreenshotFilename(shot.Name)
		subdir := filepath.Join(config.BaseDir, shot.Name)

		err := os.MkdirAll(subdir, 0774)
		if err != nil {
			log.Fatalf("unable to create directory %s", subdir)
		}

		fullpath := filepath.Join(subdir, filename)
		// capture screenshot of an element
		var buf []byte

		// capture entire browser viewport, returning png with quality=90
		if err := chromedp.Run(ctx, fullScreenshot(shot.URL, 90, &buf)); err != nil {

			log.Fatal(err)
		}
		if err := ioutil.WriteFile(fullpath, buf, 0o644); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("created screenshot %s", fullpath)
		}
	}
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

type ConfigObj struct {
	BaseDir string `json:"screenshot_base_dir"`
	Shots   []ShotObj
}
type ShotObj struct {
	Name string
	URL  string `json:"url"`
}

func loadJSONConfig(filepath string) *ConfigObj {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open config %s: %v", filepath, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("unable to parse bytevalue of %s: %v", filepath, err)
	}
	var configObj ConfigObj
	json.Unmarshal(byteValue, &configObj)

	fi, err := os.Stat(configObj.BaseDir)
	if err != nil {
		log.Fatalf("unable to stat the base directory %s", configObj.BaseDir)
	}
	if !fi.IsDir() {
		log.Fatalf("base directory is not a directory %s", configObj.BaseDir)
	}
	return &configObj
}
