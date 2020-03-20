package msdl

import (
    "fmt"
    "net/http"
    "path"
    "strings"

    "github.com/briandowns/spinner"
    "github.com/PuerkitoBio/goquery"
)

func GetPages(id string, secret string, dir string, s *spinner.Spinner) (int, error) {
	current_page := 0
	x := id[len(id) - 1:]
	y := id[len(id) - 2:len(id) - 1]
	z := id[len(id) - 3:len(id) - 2]

	for {
		current_url := fmt.Sprintf("https://musescore.com/static/musescore/scoredata/gen/%s/%s/%s/%s/%s/score_%d.svg", x, y, z, id, secret, current_page)
		resp, err := http.Head(current_url)
		if err != nil {
		    return current_page, err
		}

		if resp.StatusCode != 200 {
			break
		} else {
            s.Suffix = fmt.Sprintf(" Downloading score (page %d)", current_page + 1)
			err = GeneratePDF(current_url, fmt.Sprintf(path.Join(dir, "score_%d.pdf"), current_page))
			if err != nil {
				return current_page, err
			}

			current_page++
		}
	}
	return current_page, nil
}

func GetScoreIDSecret(url string) (string, string) {
	// Make HTTP request
    response, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()

    // Read response data in to memory
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        panic(err)
    }

    //Process meta elements
    // Find and print image URLs

    var id string
    var secret string

    document.Find("meta").Each(func(index int, element *goquery.Selection) {
        property, exists := element.Attr("property")
        if exists && property == "og:image" {
            content, exists := element.Attr("content")
            if exists {
            	urlparts := strings.Split(content, "/")
            	id = urlparts[10]
            	secret = urlparts[11]
            }
        }
    })

    return id, secret
}

func GetScoreTitle(url string) string {
        // Make HTTP request
        response, err := http.Get(url)
        if err != nil {
            panic(err)
        }
        defer response.Body.Close()

        // Read response data in to memory
        document, err := goquery.NewDocumentFromReader(response.Body)
        if err != nil {
            panic(err)
        }

        //Process meta elements
        var title string

        document.Find("meta").Each(func(index int, element *goquery.Selection) {
            property, exists := element.Attr("property")
            if exists && property == "og:title" {
                content, exists := element.Attr("content")
                if exists {
                    title = content
                }
            }
        })

        return title
}
