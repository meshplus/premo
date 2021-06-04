package interface_tester

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpGet(url string) ([]byte, error) {
	/* #nosec */
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func httpPost(url string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)

	/* #nosec */
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return nil, err
	}
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getURL(path string) string {
	return "http://172.27.239.94:9091/v1/" + path
}
