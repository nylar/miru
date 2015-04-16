package miru

import (
	"net/url"
	"strings"
)

// ProcessURL determines whether a URL is to be enqueued or not.
func ProcessURL(link, site string) (string, error) {
	if strings.HasPrefix(link, "#") {
		return "", ErrInvalidURL
	}

	// A relative link that may start with a period, remove the period.
	if strings.HasPrefix(link, ".") {
		link = link[1:]
	}

	// Normalise paths so that we explicility handle slashes
	if strings.HasPrefix(link, "/") {
		link = link[1:]
	}

	_url, _ := url.Parse(link)

	// Check url is absolute and equal to the root site.
	if _url.IsAbs() {
		if _url.Host != site {
			return "", ErrInvalidURL
		}
	} else { // Generate a full URL
		link = "http://" + site + "/" + link
	}

	return link, nil
}
