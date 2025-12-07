package atfram

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) setHeadersForRequest(req *http.Request) {
	if req != nil {
		for k, v := range c.atconf.Headers {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) DoHttpRequest(method, url string, data []byte) (*http.Response, error) {
	// Construct Request
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}

	// set default headers
	c.setHeadersForRequest(req)

	//! log
	c.logger.Debug("REQ", method, url, fmt.Sprintf("%v", len(data)))

	// Do Request
	resp, err := c.c.Do(req)
	if err != nil {
		return resp, err
	}

	for _, v := range resp.Cookies() {
		//
		var value = v.Value
		if len(v.Value) > 16 {
			value = value[:16] + "..."
		}

		c.logger.Info("GOT COOKIE", v.Name, value)
	}

	return resp, err
}

/*

	//
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(data))

	//! log
	c.logger.Debug("RESP", string(body))
*/

// ! UNSAFE
func (c *Client) FollowHttpRedirect(link string) (*http.Response, error) {
	resp, err := c.DoHttpRequest(
		http.MethodGet,
		link,
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther:
		return c.FollowHttpRedirect(
			resp.Header.Get("Location"),
		)

	case http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return c.FollowHttpRedirect(
			resp.Header.Get("Location"),
		)
	}

	return resp, nil
}

// eig nur für debug logs
func _readToBuf(r io.ReadCloser) *bytes.Buffer {

	out, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	//? io.LimitReader()

	return bytes.NewBuffer(out)
}
