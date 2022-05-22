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

var config *ConfigObj

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("you must include the config file path")
	}

	config = loadJSONConfig(os.Args[1])
	// create context

	for _, shot := range config.Shots {

		ctx, cancel := chromedp.NewContext(
			context.Background(),
			// chromedp.WithDebugf(log.Printf),
		)

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
		if err := chromedp.Run(ctx, fullScreenshot(&shot, &buf)); err != nil {

			log.Fatal(err)
		}
		if err := ioutil.WriteFile(fullpath, buf, 0o644); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("created screenshot %s", fullpath)
		}
		cancel()
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
func fullScreenshot(shot *ShotObj, res *[]byte) chromedp.Tasks {
	//
	scriptsBytes, err := ioutil.ReadFile(filepath.Join(config.ScriptsDir, shot.Script))
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	code := string(scriptsBytes)
	return chromedp.Tasks{
		chromedp.EmulateReset(),
		chromedp.ResetViewport(),
		chromedp.EmulateViewport(shot.Width, shot.Height),
		chromedp.Navigate(shot.URL),
		chromedp.Sleep(time.Second * time.Duration(shot.Sleep)),
		chromedp.Evaluate(code, nil),
		chromedp.FullScreenshot(res, shot.Quality),
	}
}

func generateScreenshotFilename(prefix string) string {
	dt := time.Now()
	timestamp := dt.Format("D01_02_2006-T15_04_05")
	return fmt.Sprintf("%s-%s.png", prefix, timestamp)
}

type ConfigObj struct {
	BaseDir    string `json:"screenshots_base_dir"`
	ScriptsDir string `json:"scripts_base_dir"`
	Shots      []ShotObj
}
type ShotObj struct {
	Name    string
	URL     string `json:"url"`
	Script  string
	Sleep   int64
	Quality int
	Width   int64
	Height  int64
}

func loadJSONConfig(configpath string) *ConfigObj {
	jsonFile, err := os.Open(configpath)
	if err != nil {
		log.Fatalf("failed to open config %s: %v", configpath, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("unable to parse bytevalue of %s: %v", configpath, err)
	}
	var configObj ConfigObj
	json.Unmarshal(byteValue, &configObj)

	// Verify dirs exist
	existsFatal(configObj.BaseDir)
	existsFatal(configObj.ScriptsDir)

	// Verify script files exist
	for _, shot := range configObj.Shots {
		existsFatal(filepath.Join(configObj.ScriptsDir, shot.Script))
	}
	return &configObj
}

// fatality if dir or file doesn't exist
func existsFatal(path string) {
	ok, err := exists(path)
	if !ok && err == nil {
		log.Fatalf("file or directory doesn't exist: %s", path)
	}
	if !ok && err != nil {
		log.Fatalf("error stating file or directory %s : %v", path, err)
	}
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
