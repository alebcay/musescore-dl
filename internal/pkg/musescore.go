package msdl

import (
    "context"
    "fmt"
    "net/http"
    "path"
    "regexp"
    "strconv"
    "strings"

    "github.com/briandowns/spinner"
    "github.com/chromedp/chromedp"
    "github.com/PuerkitoBio/goquery"
)

func GetNumberOfPages(url string) (int, error) {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var res string
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.InnerHTML(`:root`, &res, chromedp.NodeVisible),
    )
    if err != nil {
        panic(err)
    }

    re := regexp.MustCompile(`– (\d+) of (\d+) pages"\ssrc="`)
    number, err := strconv.Atoi(strings.Fields(strings.TrimSpace(re.FindString(res)))[3])

    return number, err
}

func GetPages(id string, secret string, dir string, s *spinner.Spinner, pages int) error {
	x := id[len(id) - 1:]
	y := id[len(id) - 2:len(id) - 1]
	z := id[len(id) - 3:len(id) - 2]

	for current_page := 0; current_page < pages; current_page++ {
		current_url := fmt.Sprintf("https://musescore.com/static/musescore/scoredata/gen/%s/%s/%s/%s/%s/score_%d.svg", x, y, z, id, secret, current_page)
		resp, err := http.Get(current_url)
		if err != nil {
		    return err
		}

		if resp.StatusCode != 200 {
			break
		} else {
            // Read response data in to memory
            document, err := goquery.NewDocumentFromReader(resp.Body)
            if err != nil {
                return err
            }

            var width float64
            var height float64

            document.Find("svg").Each(func(index int, element *goquery.Selection) {
                width_element, exists := element.Attr("width")
                if exists {
                    width, err = strconv.ParseFloat(width_element, 64)
                    if err != nil {
                        panic(err)
                    }
                }

                height_element, exists := element.Attr("height")
                if exists {
                    height, err = strconv.ParseFloat(height_element, 64)
                    if err != nil {
                        panic(err)
                    }
                }
            })

            s.Suffix = fmt.Sprintf(" Downloading score (page %d of %d)", current_page + 1, pages)
			err = GeneratePDF(current_url, fmt.Sprintf(path.Join(dir, "score_%d.pdf"), current_page), width, height)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
