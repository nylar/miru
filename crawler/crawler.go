package crawler

import (
	"io/ioutil"
	"net/http"
)

var UserAgent = "Miru/1.0 (+http://www.miru.nylar.io)"

func getDocument(url string) ([]byte, error) {
	client := &http.Client{}
	data := []byte{}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", UserAgent)

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	data, _ = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	return data, nil
}
