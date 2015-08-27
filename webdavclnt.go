package webdavclnt

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type WebDavClient struct {
	Host string
	Port int
	Login string
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
	if clnt.Port > 0  {
		connectionString += ":" + strconv.Itoa(clnt.Port)
	}

	return connectionString
}

func (clnt *WebDavClient) Get(uri string) error {
	return nil
}

func (clnt *WebDavClient) Upload(uri string, data io.Reader) error {

    req, err := http.NewRequest("PUT", clnt.buildConnectionString() + uri, data)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")
	req.SetBasicAuth(clnt.Login, clnt.Password)

    httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

	return nil
}

func (clnt *WebDavClient) Delete(uri string) error {
	return nil
}
