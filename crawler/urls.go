package crawler

import (
	"net/url"
	"strings"
)

func ProcessURL(link, site string) (string, error) {
	if strings.HasPrefix(link, "#") {
		return "", InvalidUrlError
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
			return "", InvalidUrlError
		}
	} else { // Generate a full URL
		link = "http://" + site + "/" + link
	}

	return link, nil
}
