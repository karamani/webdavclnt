package webdavclnt

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type WebDavClient struct {
	Host     string
	Port     int
	Login    string
	Password string
}

func NewClient(host string) *WebDavClient {
	return &WebDavClient{
		Host:     host,
		Port:     0,
		Login:    "",
		Password: "",
	}
}

func (clnt *WebDavClient) buildConnectionString() string {

	var connectionString string

	connectionString = clnt.Host
	if !strings.Contains(clnt.Host, "http://") && !strings.Contains(clnt.Host, "https://") {
		connectionString = "http://" + connectionString
	}
	if clnt.Port > 0 {
		connectionString += ":" + strconv.Itoa(clnt.Port)
	}

	return connectionString
}

func (clnt *WebDavClient) buildRequest(method, uri string, data io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, clnt.buildConnectionString()+uri, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	if len(clnt.Login) > 0 {
		req.SetBasicAuth(clnt.Login, clnt.Password)
	}

	return req, nil
}

func (clnt *WebDavClient) Get(uri string) ([]byte, error) {

	req, err := clnt.buildRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (clnt *WebDavClient) Put(uri string, data io.Reader) error {

	req, err := clnt.buildRequest("PUT", uri, data)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (clnt *WebDavClient) Delete(uri string) error {

	req, err := clnt.buildRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
