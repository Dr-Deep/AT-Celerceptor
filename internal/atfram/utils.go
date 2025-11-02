package atfram

import (
	"net/http"
)

func (c *Client) setHeadersForRequest(req *http.Request) {
	if req != nil {
		for k, v := range c.atconf.Headers {
			req.Header.Set(k, v)
		}
	}
}

/*
func readToBuf(r io.ReadCloser) (*bytes.Buffer, error) {
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	//? io.LimitReader()

	return bytes.NewBuffer(out), r.Close()
}
*/
