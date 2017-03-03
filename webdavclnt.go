//
// Package webdavclnt contains WebDav Http Client.
//
// Author: Yuri Y. Karamani <y.karamani@gmail.com>
//
package webdavclnt

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Properties map[string]string

type WebDavClient struct {
	Host      string
	Port      int
	Login     string
	Password  string
	DefFolder string
}

// NewClient creates a pointer to an instance of WebDavClient.
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

	connectionString := clnt.Host

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

// SetPort sets the value of a field Port,
// returns a poiner to an instance WebDavClient.
func (clnt *WebDavClient) SetPort(port int) *WebDavClient {
	clnt.Port = port
	return clnt
}

// SetLogin sets the value of a field Login,
// returns a poiner to an instance WebDavClient.
func (clnt *WebDavClient) SetLogin(login string) *WebDavClient {
	clnt.Login = login
	return clnt
}

// SetPassword sets the value of a field Password,
// returns a poiner to an instance WebDavClient.
func (clnt *WebDavClient) SetPassword(password string) *WebDavClient {
	clnt.Password = password
	return clnt
}

// SetDefFolder sets the value of a field DefFolder,
// returns a poiner to an instance WebDavClient.
func (clnt *WebDavClient) SetDefFolder(defFolder string) *WebDavClient {
	clnt.DefFolder = defFolder
	return clnt
}

// statusIsValid validates response status.
// Valid statuses: 2xx or 3xx
func (clnt *WebDavClient) statusIsValid(status int) bool {
	return status >= http.StatusOK && status < http.StatusBadRequest
}

// Get gets a file from uri.
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

	if !clnt.statusIsValid(resp.StatusCode) {
		return nil, errors.New("Error: " + resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// Put uploads a file into WebDav storage to the specified uri.
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

// Delete removes a file from WebDav storage to the specified uri.
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

// MkCol makes new directory (collection).
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

// Copy copies a file from uri to destURI.
func (clnt *WebDavClient) Copy(uri, destURI string) error {

	req, err := clnt.buildRequest("COPY", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Destination", clnt.buildConnectionString() + clnt.DefFolder + destURI)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Move moves a file from uri to destURI.
func (clnt *WebDavClient) Move(uri, destURI string) error {

	req, err := clnt.buildRequest("MOVE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Destination", clnt.buildConnectionString() + clnt.DefFolder + destURI)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (clnt *WebDavClient) getProps(uri, propfind string) (map[string]Properties, error) {

	body := bytes.NewBufferString(
		fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" ?><propfind xmlns="DAV:">%s</propfind>`, propfind))

	req, err := clnt.buildRequest("PROPFIND", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Depth", "1")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !clnt.statusIsValid(resp.StatusCode) {
		return nil, errors.New("Error: " + resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	obj := Multistatus{}
	err = xml.Unmarshal(contents, &obj)
	if err != nil {
		return nil, err
	}

	if obj.Responses == nil || len(obj.Responses) == 0 {
		return nil, errors.New("Unknown xml schema")
	}

	res := make(map[string]Properties)

	for _, respTag := range obj.Responses {
		if respTag.Propstat == nil || respTag.Propstat.Prop == nil ||
			respTag.Propstat.Prop.PropList == nil {

			return nil, errors.New("Unknown xml schema")
		}

		p := make(Properties)
		for _, prop := range respTag.Propstat.Prop.PropList {
			p[prop.XMLName.Local] = prop.Value
		}

		reskey := respTag.Href
		if len(clnt.DefFolder) > 0 && strings.Index(respTag.Href, clnt.DefFolder) == 0 {
			reskey = strings.Replace(reskey, clnt.DefFolder, "", 1)
		}
		res[reskey] = p
	}

	return res, nil
}

// PropFind gets specific properties of file at the specified uri.
func (clnt *WebDavClient) PropFind(uri string, props ...string) (map[string]Properties, error) {

	propstr := "<prop>"
	for _, eachProp := range props {
		propstr += fmt.Sprintf("<%s/>", eachProp)
	}
	propstr += "</prop>"

	return clnt.getProps(uri, propstr)
}

// AllPropFind gets all properties of file at the specified uri.
func (clnt *WebDavClient) AllPropFind(uri string) (map[string]Properties, error) {
	return clnt.getProps(uri, "<allprop/>")
}

// PropNameFind gets properties names of file at the specified uri.
func (clnt *WebDavClient) PropNameFind(uri string) (map[string][]string, error) {

	props, err := clnt.getProps(uri, "<propname/>")
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	for respkey, resp := range props {
		for key := range resp {
			res[respkey] = append(res[respkey], key)
		}
	}

	return res, nil
}

// Exists checks if the file or collection exists.
func (clnt *WebDavClient) Exists(uri string) (bool, error) {

	req, err := clnt.buildRequest("HEAD", uri, nil)
	if err != nil {
		return false, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
