// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package prpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/luci/luci-go/common/clock"
)

// This file implements decoding of HTTP requests to RPC parameters.

const (
	// headerSuffixBinary is a suffix of an HTTP header that specifies that
	// the header value is encoded in std base64.
	// After decoding, a handler must process the header without the suffix.
	headerSuffixBinary = "-Bin"
	headerContentType  = "Content-Type"
)

// requestFormat converts Content-Type header value from a request to a format.
// Can return only formatBinary, formatJSONPB or formatText.
// In case of an error, format is undefined.
func requestFormat(contentType string) (format, error) {
	if contentType == "" {
		return formatBinary, nil
	}
	mediaType, mediaTypeParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		return formatBinary, err
	}
	f, err := parseFormat(mediaType, mediaTypeParams)
	if err != nil {
		return f, err
	}

	switch f {
	case formatUnrecognized:
		return f, fmt.Errorf("%q is not supported", contentType)

	case formatUnspecified:
		return formatBinary, nil

	case formatBinary, formatJSONPB, formatText:
		return f, nil

	default:
		panic("cannot happen")
	}
}

// readMessage decodes a protobuf message from an HTTP request.
// Does not close the request body.
func readMessage(r *http.Request, msg proto.Message) *httpError {
	if msg == nil {
		panicf("cannot decode to nil")
	}
	format, err := requestFormat(r.Header.Get(headerContentType))
	if err != nil {
		// Spec: http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.16
		return errorf(http.StatusUnsupportedMediaType, "Content-Type header: %s", err)
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errorf(http.StatusBadRequest, "could not read body: %s", err)
	}
	if len(buf) == 0 {
		// no body, no message
		return nil
	}

	switch format {
	// Do not redefine "err" below.

	case formatBinary:
		err = proto.Unmarshal(buf, msg)

	case formatJSONPB:
		err = jsonpb.Unmarshal(bytes.NewBuffer(buf), msg)

	case formatText:
		err = proto.UnmarshalText(string(buf), msg)

	default:
		panicf("cannot happen")
	}
	if err != nil {
		return errorf(http.StatusBadRequest, "could not decode body: %s", err)
	}
	return nil
}

// parseHeader parses HTTP headers and derives a new context.
// Supports headerTimeout.
// Ignores "Accept" and "Content-Type" headers.
//
// If there are unrecognized HTTP headers, with or without headerSuffixBinary,
// they are added to a metadata.MD and a new context is derived.
// If c already has metadata, the latter is copied.
//
// In case of an error, returns c unmodified.
func parseHeader(c context.Context, header http.Header) (context.Context, error) {
	origC := c

	md, ok := metadata.FromContext(c)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.MD{}
	}

	addedMeta := false
	for name, values := range header {
		if len(values) == 0 {
			continue
		}
		name = http.CanonicalHeaderKey(name)
		switch name {

		case headerTimeout:
			// Decode only first value, ignore the rest
			// to be consistent with http.Header.Get.
			timeout, err := decodeTimeout(values[0])
			if err != nil {
				return origC, fmt.Errorf("%s header: %s", headerTimeout, err)
			}
			c, _ = clock.WithTimeout(c, timeout)

		case headerAccept, headerContentType:
		// readMessage and writeMessage handle these headers.

		default:
			addedMeta = true
			if !strings.HasSuffix(name, headerSuffixBinary) {
				md[name] = append(md[name], values...)
				break // switch name
			}
			trimmedName := strings.TrimSuffix(name, headerSuffixBinary)
			for _, v := range values {
				decoded, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					return origC, fmt.Errorf("%s header: %s", name, err)
				}
				md[trimmedName] = append(md[trimmedName], string(decoded))
			}
		}
	}
	if addedMeta {
		c = metadata.NewContext(c, md)
	}
	return c, nil
}