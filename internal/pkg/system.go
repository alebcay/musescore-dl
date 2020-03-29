package msdl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/mholt/archiver/v3"
)

var ctx context.Context
var once sync.Once

func SetupChrome(dir string) error {
	zipfile, err := assets.Open("chrome.zip")
	if err != nil {
		return err
	}
	defer zipfile.Close()

	dest_copy, err := os.OpenFile(path.Join(dir, "chrome.zip"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer dest_copy.Close()

	_, err = io.Copy(dest_copy, zipfile)
	if err != nil {
		return err
	}

	archiver.Unarchive(path.Join(dir, "chrome.zip"), dir)

	switch runtime.GOOS {
	case "darwin":
		final_path := path.Join(dir, "chrome-mac", "Chromium.app", "Contents", "MacOS", "Chromium")
		err = WriteChromeShimScript(path.Join(dir, "google-chrome"), final_path)
		os.Setenv("PATH", dir+":"+os.Getenv("PATH"))
	case "linux":
		final_path := path.Join(dir, "chrome-linux", "chrome")
		err = WriteChromeShimScript(path.Join(dir, "google-chrome"), final_path)
		os.Setenv("PATH", dir+":"+os.Getenv("PATH"))
	case "windows":
		base_dir, _ := path.Split(path.Join(dir, "chrome-win"))
		os.Setenv("PATH", base_dir+":"+os.Getenv("PATH"))
	default:
		return errors.New(fmt.Sprintf("unsupported platform %s", runtime.GOOS))
	}

	return err
}

func WriteChromeShimScript(file string, exec_path string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("#!/bin/sh\n%s $@\n", exec_path))
	if err != nil {
		return err
	}
	os.Chmod(file, 0777)
	return err
}

func GetChromeContext() context.Context {
    once.Do(func() {
        ctx, _ = chromedp.NewContext(context.Background())
    })
    return ctx
}
