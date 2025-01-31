package menandmice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	ResourceNotFound = 16544
)

type Mmclient struct{ resty.Client }

// Cfg config to construct client
type Cfg struct {
	MMEndpoint string
	MMUsername string
	MMPassword string
	TLSVerify  bool
	Timeout    int
}

func init() {
	// Remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// ClientInit establishes default settings on the REST client
func ClientInit(c *Cfg) (*Mmclient, error) {
	client := Mmclient{Client: *resty.New()}

	if c.MMEndpoint == "" {
		return nil, errors.New("REST API endpoint must be configured")
		//TODO check if it resolaves
	}

	if match, _ := regexp.MatchString("^(http|https)://", c.MMEndpoint); !match {

		return nil, fmt.Errorf("REST API endpoint: %s must start with \"http://\" or \"https://\"", c.MMEndpoint)
	}

	if c.MMUsername == "" {
		return nil, errors.New("Invalid username")
	}
	if c.MMPassword == "" {
		return nil, errors.New("Invalid password")
	}

	if c.TLSVerify == false {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	}

	client.SetBasicAuth(c.MMUsername, c.MMPassword)
	client.SetHeader("Content-Type", "application/json")
	client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	client.SetHostURL(c.MMEndpoint + "/mmws/api")

	// TODO check if this works well with Round Robin DNS
	// TODO is this needed
	client.SetRetryCount(5)
	client.SetRetryWaitTime(1 * time.Second)
	client.AddRetryCondition(func(r *resty.Response, e error) bool {
		// also retry  on server errors
		return r.StatusCode() >= 500 && r.StatusCode() < 600
	})

	// Test if we can make a connection

	// TODO use request that need authentication
	_, err := client.R().Get("")
	if err != nil {
		return nil, fmt.Errorf("Could not connect with endpoint: %s\n\t%s", c.MMEndpoint, err)
	}

	return &client, err
}

type DeleteRequest struct {
	SaveComment  string `json:"saveComment"`
	ForceRemoval bool   `json:"forceRemoval"`
	ObjType      string `json:"objType,omitempty"`
}

func deleteRequest(objType string) DeleteRequest {
	return DeleteRequest{
		ForceRemoval: true,
		SaveComment:  "deleted by terraform",
		ObjType:      objType,
	}
}

type RefResponse struct {
	Result struct {
		Ref string `json:"ref"`
	} `json:"result"`
}
type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type RequestError struct {
	Method     string
	URL        string
	HTTPCode   int
	StatusCode int
	ErrMessage string
}

func (r *RequestError) Error() string {

	var operation string
	switch r.Method {
	case "GET":
		operation = "reading"

	case "PUT":
		operation = "updating"

	case "POST":
		operation = "creating"

	case "DELETE":
		operation = "deleting"
	default:
		operation = "accesing"
	}

	url, err := url.Parse(r.URL)
	if err != nil {
		log.Fatal(err)
	}
	resource := url.RequestURI()

	return fmt.Sprintf("Failed with %v %v\n\tHTTP code:%v: %v", operation, resource, r.HTTPCode, r.ErrMessage)
}

func ResponseError(response *resty.Response, errorResponse ErrorResponse) error {

	if !response.IsSuccess() {
		return &RequestError{

			Method:     response.Request.Method,
			URL:        response.Request.URL,
			HTTPCode:   response.StatusCode(),
			StatusCode: errorResponse.Error.Code,
			ErrMessage: errorResponse.Error.Message,
		}
	}
	return nil
}

func (c *Mmclient) Get(result interface{}, path string, query map[string]interface{}, filter map[string]string) error {

	//TODO better error message
	var errorResponse ErrorResponse
	var querystring string

	request := c.R().
		SetError(&errorResponse).
		SetResult(&result)

	if query != nil {
		for key, val := range query {

			request = request.SetQueryParam(key, fmt.Sprintf("%v", val))
		}
	}
	if filter != nil {

		conditions := make([]string, 0, len(filter))
		for key, val := range filter {
			conditions = append(conditions, fmt.Sprintf("%s=%s", key, val))
		}
		querystring = strings.Join(conditions, "&")
		request = request.SetQueryParam("filter", querystring)
	}

	r, err := request.Get(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

func (c *Mmclient) Post(data interface{}, result interface{}, path string) error {

	//TODO better error message
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		SetResult(&result).
		Post(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

func (c *Mmclient) Delete(data interface{}, path string) error {

	var err error
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Delete(path)

	if err != nil {
		return err
	}

	return ResponseError(r, errorResponse)
}

func (c *Mmclient) Put(data interface{}, path string) error {
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(data).
		SetError(&errorResponse).
		Put(path)

	if err != nil {
		return err
	}
	return ResponseError(r, errorResponse)
}
