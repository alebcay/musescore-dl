package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/shurcooL/vfsgen"
)

func main() {
	current_dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	assets_path := path.Join(current_dir, "build", "deps")
	os.MkdirAll(assets_path, 0755)
	err = GetChrome(assets_path)
	if err != nil {
		panic(err)
	}

	generated_file_path, err := filepath.Abs(path.Join(current_dir, "internal", "pkg", "chrome.go"))
	if err != nil {
		panic(err)
	}

	fs := http.Dir(assets_path)

	err = vfsgen.Generate(fs, vfsgen.Options{
		Filename:    generated_file_path,
		PackageName: "msdl",
	})
	if err != nil {
		panic(err)
	}
}

func GetChrome(dir string) error {
	var platform string
	var filename string

	switch runtime.GOOS {
	case "darwin":
		platform = "Mac"
		filename = "mac"
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			platform = "Linux_x64"
			filename = "linux"
		case "386":
			platform = "Linux"
			filename = "linux"
		default:
			return errors.New(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			platform = "Win_x64"
			filename = "win"
		case "386":
			platform = "Win"
			filename = "win"
		default:
			return errors.New(fmt.Sprintf("unsupported architecture %s", runtime.GOARCH))
		}
	default:
		return errors.New(fmt.Sprintf("unsupported platform %s", runtime.GOOS))
	}

	last_change_url := fmt.Sprintf("https://commondatastorage.googleapis.com/chromium-browser-snapshots/%s/LAST_CHANGE", platform)
	resp, err := http.Get(last_change_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var last_change_string string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		last_change_string = string(bodyBytes)
	} else {
		return errors.New("unable to determine latest chromium version")
	}

	resp.Body.Close()

	archive_url := fmt.Sprintf("https://commondatastorage.googleapis.com/chromium-browser-snapshots/%s/%s/chrome-%s.zip", platform, last_change_string, filename)
	resp, err = http.Get(archive_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	zip_path := path.Join(dir, "chrome.zip")
	out, err := os.Create(zip_path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}
