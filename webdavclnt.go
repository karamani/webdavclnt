//
// WebDav Http Client
//
// Author: Yuri Y. Karamani <y.karamani@gmail.com>
//
package webdavclnt

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type WebDavClient struct {
	Host      string
	Port      int
	Login     string
	Password  string
	DefFolder string
}

//
// WebDav Client constructor
//
func NewClient(host string) *WebDavClient {
	return &WebDavClient{
		Host:      host,
		Port:      0,
		Login:     "",
		Password:  "",
		DefFolder: "",
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

	req, err := http.NewRequest(method, clnt.buildConnectionString()+clnt.DefFolder+uri, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	if len(clnt.Login) > 0 {
		req.SetBasicAuth(clnt.Login, clnt.Password)
	}

	return req, nil
}

func (clnt *WebDavClient) SetPort(port int) *WebDavClient {
	clnt.Port = port
	return clnt
}

func (clnt *WebDavClient) SetLogin(login string) *WebDavClient {
	clnt.Login = login
	return clnt
}

func (clnt *WebDavClient) SetPassword(password string) *WebDavClient {
	clnt.Password = password
	return clnt
}

func (clnt *WebDavClient) SetDefFolder(defFolder string) *WebDavClient {
	clnt.DefFolder = defFolder
	return clnt
}

//
// Get file from WebDav Storage
//
func (clnt *WebDavClient) Get(uri string) ([]byte, error) {

	req, err := clnt.buildRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
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

//
// Upload file into WebDav Storage
//destUri
func (clnt *WebDavClient) Put(uri string, data io.Reader) error {

	req, err := clnt.buildRequest("PUT", uri, data)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//
// Delete file from WebDav Storage
//
func (clnt *WebDavClient) Delete(uri string) error {

	req, err := clnt.buildRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//
// Make new directory (collection)
//
func (clnt *WebDavClient) MkCol(uri string) error {

	req, err := clnt.buildRequest("MKCOL", uri, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//
// Copy file
//
func (clnt *WebDavClient) Copy(uri, destUri string) error {

	req, err := clnt.buildRequest("COPY", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Destination", clnt.buildConnectionString()+clnt.DefFolder+destUri)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//
// Move file
//
func (clnt *WebDavClient) Move(uri, destUri string) error {

	req, err := clnt.buildRequest("MOVE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Destination", clnt.buildConnectionString()+clnt.DefFolder+destUri)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
