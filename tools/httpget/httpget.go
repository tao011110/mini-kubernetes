package httpget

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

/*
	HOW TO USE:
		err:=httpclient.Post("http://localhost:8080/...").
            ContentType("application/json").
            Body(`{"name":"lewiskong"}`).
            AddHeader("Referer","localhost").
            AddCookie("uin","guest").
            GetString(&str).
            GetJson(&obj). //when the result is json
            Execute()
*/

type HttpMethod string

const (
	httpGet    HttpMethod = "GET"
	httpPost   HttpMethod = "POST"
	httpPut    HttpMethod = "PUT"
	httpDelete HttpMethod = "DELETE"
)

type HttpClient struct {
	Error    error
	Response *http.Response

	tasks   map[*interface{}]string
	method  HttpMethod
	request *http.Request
	routeip string
}

func Get(rawurl string) *HttpClient {
	client := CreateDefault()
	client.method = httpGet

	client.request, client.Error = http.NewRequest("GET", "", nil)
	if client.Error != nil {
		return client
	}

	rawurl, err := HandleURL(client, rawurl)
	if err != nil {
		client.Error = err
	}

	return client.handle(rawurl)
}

func Post(rawurl string) *HttpClient {
	client := CreateDefault()
	client.method = httpPost
	client.request, client.Error = http.NewRequest("POST", "", nil)
	if client.Error != nil {
		return client
	}

	rawurl, err := HandleURL(client, rawurl)
	if err != nil {
		client.Error = err
	}

	return client.handle(rawurl)
}

func Put(rawurl string) *HttpClient {
	client := CreateDefault()
	client.method = httpPut
	client.request, client.Error = http.NewRequest("PUT", "", nil)
	if client.Error != nil {
		return client
	}

	rawurl, err := HandleURL(client, rawurl)
	if err != nil {
		client.Error = err
	}

	return client.handle(rawurl)
}

func DELETE(rawurl string) *HttpClient {
	client := CreateDefault()
	client.method = httpDelete

	client.request, client.Error = http.NewRequest("DELETE", "", nil)
	if client.Error != nil {
		return client
	}

	rawurl, err := HandleURL(client, rawurl)
	if err != nil {
		client.Error = err
	}

	return client.handle(rawurl)
}

func CreateDefault() *HttpClient {
	client := new(HttpClient)
	client.tasks = map[*interface{}]string{}
	return client
}

func (client *HttpClient) handle(rawurl string) *HttpClient {

	if client.Error != nil {
		return client
	}

	client.request.URL, client.Error = url.ParseRequestURI(rawurl)
	if client.Error != nil {
		return client
	}
	// fmt.Println(client.request.URL)
	return client
}

// UseHTTPS ...
// Whether to use https protocol
func (client *HttpClient) UseHTTPS(flag bool) *HttpClient {
	if flag {
		client.request.URL.Scheme = "https"
	} else {
		client.request.URL.Scheme = "http"
	}
	return client
}

// ContentType ...
// Set the content type of the request
func (client *HttpClient) ContentType(contentType string) *HttpClient {
	if client.Error != nil {
		return client
	}

	client.request.Header.Set("Content-Type", contentType)

	return client
}

// Body ...
// Set the body of the request .
// Used when the method is post .
func (client *HttpClient) Body(body io.Reader) *HttpClient {
	if client.Error != nil {
		return client
	}

	client.request.Body = ioutil.NopCloser(body)

	return client
}

// AddHeader ...
//	AddHeader to the request
func (client *HttpClient) AddHeader(key, value string) *HttpClient {
	if client.Error != nil {
		return client
	}

	client.request.Header.Add(key, value)

	return client
}

// AddCookie ...
// Add cookie to the request
func (client *HttpClient) AddCookie(name, value string) *HttpClient {
	if client.Error != nil {
		return client
	}

	ck := new(http.Cookie)
	ck.Name = name
	ck.Value = value

	client.request.AddCookie(ck)

	return client
}

func (client *HttpClient) GetString(v interface{}) *HttpClient {
	if client.Error != nil {
		return client
	}
	client.tasks[&v] = "string"
	return client
}

func (client *HttpClient) GetJson(v interface{}) *HttpClient {
	if client.Error != nil {
		return client
	}
	client.tasks[&v] = "json"
	return client
}

func (client *HttpClient) GetJsonp(v interface{}) *HttpClient {
	if client.Error != nil {
		return client
	}
	client.tasks[&v] = "jsonp"
	return client
}

func (client *HttpClient) GetJce(v interface{}) *HttpClient {
	if client.Error != nil {
		return client
	}
	client.tasks[&v] = "jce"
	return client
}

func (client *HttpClient) Execute() (err error, status string) {

	if client.Error != nil {
		return client.Error, "500"
	}

	c := http.DefaultClient
	rsp, err := c.Do(client.request)

	if err != nil {
		fmt.Println(err)
		return err, rsp.Status[:3]
	}
	content, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err, rsp.Status[:3]
	}
	//debug
	// fmt.Println(string(content[:]))
	for value, tp := range client.tasks {
		switch tp {
		case "string":
			err := setString(value, content)
			if err != nil {
				return err, rsp.Status[:3]
			}

		case "json":
			err := setJSON(value, content)
			if err != nil {
				return err, rsp.Status[:3]
			}

		case "jsonp":
			err := setJSONP(value, content)
			if err != nil {
				return err, rsp.Status[:3]
			}

		case "jce":

		default:
		}
	}
	if err != nil {
		return err, rsp.Status[:3]
	}
	//debug
	// fmt.Println(string(content[:]))
	return nil, rsp.Status[:3]
}

func setString(v *interface{}, content []byte) error {
	interf := *v
	vtype := reflect.TypeOf(interf)
	vvalue := reflect.ValueOf(interf)

	if vtype.Kind() != reflect.Ptr {
		return fmt.Errorf("error happened when parse json : the param obj %s must be pointer", vtype.Kind().String())
	}
	vvalue = reflect.Indirect(vvalue)
	vvalue.SetString(string(content[:]))
	return nil
}

// SetJSON hello
func setJSON(v *interface{}, content []byte) error {
	interf := *v
	err := json.Unmarshal(content, interf)
	return err
}

func setJSONP(v *interface{}, content []byte) error {
	interf := *v
	str := string(content[:])
	start := strings.Index(str, "{")
	end := strings.LastIndex(str, "}")
	if start < 0 || end < 0 {
		return fmt.Errorf("Parse jsonp error , wrong jsonp format : %s ", str)
	}
	err := json.Unmarshal(content[start:end+1], interf)
	return err
}

func setJce(v *interface{}) error {
	return nil
}

/*HandleURL ...
 *	 	used to convert rawurl to real url, support l5 && zkname && http && https
 *		l5 :
 *			l5://11111:22222/test
 *		zkname:
 *			zkname://test.zkname.etc/test
 */
func HandleURL(client *HttpClient, rawurl string) (string, error) {
	rawurl = strings.TrimSpace(rawurl)
	parts := strings.Split(rawurl, "://")
	if len(parts) != 2 {
		return "", errors.New("invalid url. Please start with http/https/zkname/l5")
	}
	switch parts[0] {
	case "http":
		key := strings.Split(parts[1], "/")[0]
		client.routeip = strings.Split(key, ":")[0]
		return rawurl, nil
	case "https":
		key := strings.Split(parts[1], "/")[0]
		client.routeip = strings.Split(key, ":")[0]
		return rawurl, nil
	default:
		return "", fmt.Errorf("wrong url protocol %s. Please start with http/https", parts[0])
	}
}
