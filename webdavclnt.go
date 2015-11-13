//
// WebDav Http Client
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
// Validate response status.
// Valid statuses: 2xx or 3xx
//
func (clnt *WebDavClient) statusIsValid(status int) bool {
	return status >= http.StatusOK && status < http.StatusBadRequest
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

	if !clnt.statusIsValid(resp.StatusCode) {
		return nil, errors.New("Error: " + resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

//
// Upload file into WebDav Storage
//
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
	if err != nil{
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

//
// Find properties
//
func (clnt *WebDavClient) PropFind(uri string, props ...string) (map[string]Properties, error) {

	propstr := "<prop>"
	for _, eachProp := range props {
		propstr += fmt.Sprintf("<%s/>", eachProp)
	}
	propstr += "</prop>"

	return clnt.getProps(uri, propstr)
}

//
// Get all properties
//
func (clnt *WebDavClient) AllPropFind(uri string) (map[string]Properties, error) {
	return clnt.getProps(uri, "<allprop/>")
}


//
// Get names of properties
//
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
