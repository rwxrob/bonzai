// Copyright 2022 web Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package web provides high-level functions that are called from the Go
// Bonzai branch of the same name providing universal access to the core
// functionality.
package web

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	rwxjson "github.com/rwxrob/bonzai/json"
	"gopkg.in/yaml.v3"
)

// TimeOut is a package global timeout for any of the high-level https
// query functions in this package. The default value is 60 seconds.
var TimeOut int = 60

// HTTPError is an error for anything other than an HTTP response in the
// 200-299 range including the 300 redirects (which should be handled by
// the Req.Submit successfully before returning). http.Response is
// embedded directly for details.
type HTTPError struct {
	Resp *http.Response
}

// Error fulfills the error interface.
func (e HTTPError) Error() string { return e.Resp.Status }

// ReqSyntaxError is for any error involving the incorrect
// definition of Req fields (such as including a question mark in
// the URL, etc.).
type ReqSyntaxError struct {
	Message string
}

// Error fulfills the error interface.
func (e ReqSyntaxError) Error() string { return e.Message }

// Client provides a way to change the default HTTP client for any
// further package HTTP request function calls. The Client can also be
// set in any Req by assigning the to the field of the same name. By
// default, it is set to http.DefaultClient. This is particularly useful
// when creating mockups and other testing.
var Client = http.DefaultClient

// Head contains headers to be added to a Req. Unlike the
// specification, only one header of a give name is allowed. For more
// precision the net/http library directly should be used instead.
type Head map[string]string

// Req is a human-friendly way to think of web requests. This design
// is closer a pragmatic curl requests than the canonical specification
// (unique headers, for example). The type and parameters of the web
// request and response are derived from the Req fields. The
// single-letter struct fields are terse but convenient. Use net/http
// when more precision is required.
//
// The body (B) can be one of several types that till trigger what is
// submitted as data portion of the request:
//
//	url.Values - triggers x-www-form-urlencoded
//	byte       - uuencoded binary data
//	string     - plain text
//
// Note that Req has no support for multi-part MIME. Use net/http
// directly if such is required.
//
// The data (D) field can also be any of several types that trigger how
// the received data is handled:
//
//	[]byte           - uudecoded binary
//	string           - plain text string
//	io.Writer        - keep as is
//	json.This        - unmarshaled JSON data into This
//	any              - unmarshaled JSON data
//
// Passing the query string as url.Values automatically add
// a question mark (?) followed by the URL encoded values to the end of
// the URL which may present a problem if the URL already has a query
// string. Encouraging the use of url.Values for passing the query
// string serves as a reminder that all query strings should be URL
// encoded (as is often forgotten).
type Req struct {
	U string          // base url, optional query string
	D any             // data to be populated and/or overwritten
	M string          // all caps method (default: GET)
	Q url.Values      // query string to append to URL (if none already)
	H Head            // header map, never more than one of same
	B any             // body data, url.Values will x-www-form-urlencoded
	C context.Context // trigger requests with context
	R *http.Response  // actual http.Response
}

// Submit synchronously sends the Req to server and populates the
// response from the server into D. Anything but a response in the
// 200s will result in an HTTPError. See Req for details on how
// inspection of Req will change the behavior of Submit
// automatically. It Req.C is nil a context.WithTimeout will
// be used and with the value of web.TimeOut.
func (req *Req) Submit() error {

	if req.M == "" {
		req.M = `GET`
	}
	req.M = strings.ToUpper(req.M)

	if strings.Index(req.U, "?") == 0 && req.Q != nil {
		req.U = req.U + "?" + req.Q.Encode()
	}

	var bodyReader io.Reader

	if req.H == nil {
		req.H = Head{}
	}

	var buf string

	switch v := req.B.(type) {
	case url.Values:
		buf = v.Encode()
		req.H["Content-Type"] = "application/x-www-form-urlencoded"
	case []byte:
		log.Println("planned, but unimplemented, would uuencode")
		//req.H["Content-Length"] = strconv.Itoa(len(uuencoded))
	case string:
		buf = v
	case json.Marshaler:
		byt, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf = string(byt)
	case encoding.TextMarshaler:
		byt, err := v.MarshalText()
		if err != nil {
			return err
		}
		buf = string(byt)
	case fmt.Stringer:
		buf = v.String()
	default:
		buf = fmt.Sprintf("%v", v)
	}

	bodyReader = strings.NewReader(buf)
	req.H["Content-Length"] = strconv.Itoa(len(buf))

	httpreq, err := http.NewRequest(req.M, req.U, bodyReader)
	if err != nil {
		return err
	}

	if req.H != nil {
		for k, v := range req.H {
			httpreq.Header.Add(k, v)
		}
	}

	if req.C == nil {
		dur := time.Duration(time.Second * time.Duration(TimeOut))
		ctx, cancel := context.WithTimeout(context.Background(), dur)
		defer cancel()
		httpreq = httpreq.WithContext(ctx)
	}

	res, err := Client.Do(httpreq)
	req.R = res

	if err != nil {
		return err
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return HTTPError{res}
	}

	resbytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(resbytes) == 0 {
		return nil
	}

	switch v := req.D.(type) {
	case map[string]any:
		return yaml.Unmarshal(resbytes, v)
	case *string:
		*v = string(resbytes)
	case []byte:
		log.Println("planned, but unimplemented, would uuencode")
		// v = uudecode(resbytes)
	case io.Writer:
		return nil
	case rwxjson.This:
		log.Println("rwxjson, planned, but unimplemented")
	default:
		return yaml.Unmarshal(resbytes, v)
	}

	return nil

}
