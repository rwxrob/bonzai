package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TimeOut is a package global timeout for any of the high-level https
// query functions in this package. The default value is 60 seconds.
var TimeOut int = 60

// Client provides a way to change the default HTTP client for
// any further package HTTP request function calls. By default, it is
// set to http.DefaultClient. This is particularly useful when creating
// mockups and other testing.
var Client = http.DefaultClient

// Request is a human-friendly way to think of web requests and the
// resulting JSON unmarshaled response. This design more similar to
// a pragmatic curl request than the canonical specification (unique
// headers, for example).
//
// Note that passing the query string as url.Values automatically add
// a question mark (?) followed by the URL encoded values to the end of
// the URL which may present a problem if the URL already has a query
// string. Encouraging the use of url.Values for passing the query
// string serves as a reminder that all query strings should be URL
// encoded (as is often forgotten).
type Request struct {
	Method string            // GET, POST, etc. (will upper)
	URL    string            // base url with no query string
	Query  url.Values        // query string to append to URL
	Header map[string]string // never more than one of same
	Body   url.Values        // body data, will JSON encode
	Into   any               // pointer to struct for unmarshaling
}

// Fetch passes the Request Client and unmarshals the JSON response into
// the data struct passed by pointer. Only new data will be unmarshaled
// leaving any existing data alone.
//
// All output data is expected to be JSON with the appropriate
// headers added to indicate it.
//
// If a Body is sent, it will be encoded as if submit from a POST form.
//
// Fetch observes the package global json.TimeOut.
//
// Status codes not in th 200s range will return an error with the
// status message.
//
// The http.DefaultClient is used by default but can be changed by
// setting json.Client.
func Fetch(it *Request) error {
	var err error
	var bodyreader io.Reader
	var bodylength string

	it.URL = it.URL + "?" + it.Query.Encode()
	if it.Method == "" {
		it.Method = `GET`
	}

	if it.Body != nil {
		encoded := it.Body.Encode()
		bodyreader = strings.NewReader(encoded)
		bodylength = strconv.Itoa(len(encoded))
	}

	req, err := http.NewRequest(it.Method, it.URL, bodyreader)
	if it.Body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", bodylength)
	}
	if err != nil {
		return err
	}

	if it.Header != nil {
		for k, v := range it.Header {
			req.Header.Add(k, v)
		}
	}

	dur := time.Duration(time.Second * time.Duration(TimeOut))
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	req = req.WithContext(ctx)

	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return fmt.Errorf(res.Status)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, it.Into)
}
