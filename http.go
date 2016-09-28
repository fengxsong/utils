package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var UA = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.82 Safari/537.36"

func HttpDo(c *http.Client, method, url string, header http.Header, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UA)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	if resp.StatusCode == 404 {
		err = fmt.Errorf("Not Found: %s", url)
	} else {
		err = fmt.Errorf("%s %s -> %d", method, url, resp.StatusCode)
	}
	return nil, err
}

func HttpGet(c *http.Client, url string, header http.Header) (io.ReadCloser, error) {
	return HttpDo(c, "GET", url, header, nil)
}

func HttpPost(c *http.Client, url string, header http.Header, body []byte) (io.ReadCloser, error) {
	return HttpDo(c, "POST", url, header, bytes.NewBuffer(body))
}

func HttpGetToFile(c *http.Client, url string, header http.Header, filename string) error {
	rc, err := HttpGet(c, url, header)
	if err != nil {
		return err
	}
	defer rc.Close()
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return err
}

func HttpGetBytes(c *http.Client, url string, header http.Header) ([]byte, error) {
	rc, err := HttpGet(c, url, header)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func HttpGetJSON(c *http.Client, url string, v interface{}) error {
	rc, err := HttpGet(c, url, nil)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		return fmt.Errorf("JSON SyntaxError at %s", url)
	}
	return nil
}

func HttpPostJSON(c *http.Client, url string, body, v interface{}) error {
	j, err := json.Marshal(body)
	if err != nil {
		return err
	}
	rc, err := HttpPost(c, url, http.Header{"content-type": []string{"application/json"}}, j)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		return fmt.Errorf("JSON SyntaxError at %s", url)
	}
	return nil
}

type RawFile interface {
	Name() string
	RawUrl() string
	Data() []byte
	SetData([]byte)
}

func FetchFiles(c *http.Client, files []RawFile, header http.Header) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			p, err := HttpGetBytes(c, files[i].RawUrl(), nil)
			if err != nil {
				ch <- err
				return
			}
			files[i].SetData(p)
			ch <- nil
		}(i)
	}
	for _ = range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

func FetchFilesCurl(files []RawFile, curlOptions ...string) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			stdout, _, err := ExecCmd("curl", append(curlOptions, files[i].RawUrl())...)
			if err != nil {
				ch <- err
				return
			}
			files[i].SetData([]byte(stdout))
			ch <- nil
		}(i)
	}
	for _ = range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}
