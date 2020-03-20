package msdl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"path"

	"github.com/mholt/archiver/v3"
)

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

func GetChrome(dir string) (string, error) {
	var platform string
	var filename string
	var final_path string

	switch runtime.GOOS {
	case "darwin":
		final_path = path.Join(dir, "chrome-mac", "Chromium.app", "Contents", "MacOS", "Chromium")
		platform = "Mac"
		filename = "mac"
	case "linux":
		final_path = path.Join(dir, "chrome-linux", "chrome")

		switch runtime.GOARCH {
		case "amd64":
			platform = "Linux_x64"
			filename = "linux"
		case "386":
			platform = "Linux"
			filename = "linux"
		default:
			return "", errors.New(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
		}
	case "windows":
		final_path = path.Join(dir, "chrome-win", "chrome.exe")

		switch runtime.GOARCH {
		case "amd64":
			platform = "Win_x64"
			filename = "win"
		case "386":
			platform = "Win"
			filename = "win"
		default:
			return "", errors.New(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
		}
	default:
		return "", errors.New(fmt.Sprintf("unsupported platform %s", runtime.GOOS))
	}

	last_change_url := fmt.Sprintf("https://commondatastorage.googleapis.com/chromium-browser-snapshots/%s/LAST_CHANGE", platform)
	resp, err := http.Get(last_change_url)
	if err != nil {
	    return "", err
	}
	defer resp.Body.Close()

	var last_change_string string
	if resp.StatusCode == http.StatusOK {
	    bodyBytes, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
	        return "", err
	    }
	    last_change_string = string(bodyBytes)
	} else {
		return "", errors.New("unable to determine latest chromium version")
	}

	resp.Body.Close()

	archive_url := fmt.Sprintf("https://commondatastorage.googleapis.com/chromium-browser-snapshots/%s/%s/chrome-%s.zip", platform, last_change_string, filename)
	resp, err = http.Get(archive_url)
	if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    zip_path := path.Join(dir, "chrome.zip")
    out, err := os.Create(zip_path)
    if err != nil {
        return "", err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
    	return "", err
    }

    out.Close()

    err = archiver.Unarchive(zip_path, dir)

    return final_path, err
}
