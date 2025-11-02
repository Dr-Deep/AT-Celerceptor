package atfram

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Dr-Deep/logging-go"
)

// URLS:
//- https://login.alditalk-kundenbetreuung.de/
//- https://alditalk-kundenbetreuung.de/

// Resource-Verben: create, read, update, delete, patch, _action=script, query

/*
wir können:
- OIDC Authorization Code (empfohlen für Web Apps) oder
- REST Authenticate + Session Cookie (nützlich für CLI / Legacy integrations / automatisierte logins)

* Flow:
POST: /signin/json/realms/root/realms/alditalk/authenticate
=> authID bzw JWT Token?
=> callbacks
=> stage?
=> header?
*/

// ? 	url.URL
type Client struct {

	// authID?
	// clientID?
	// customerID?
	// etc

	c      *http.Client
	logger *logging.Logger
	atconf *AldiTalkConfig
}

/*
? mit json body bzw filled callbacks
* /authenticate
"authId": "",
"callbacks": [],
"stage": "",
"header": ""
*/
func (c *Client) authenticate() (*Response, error) {
	var jsonBody []byte

	// HTTP POST
	req, err := http.NewRequest(
		"POST",
		ALDITALK_AUTHENTICATE_EP,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	c.setHeadersForRequest(req)

	// Do Request
	rawresp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	// valid Resp?
	if rawresp.StatusCode != http.StatusOK {
		return nil, ErrAldiTalkClientInvalidStatusCode
	}

	var resp Response
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

/*
* /sessions
 */
func (c *Client) sus()

/*
	c.doc = resp
	c.httpStatusCode = resp.StatusCode
*/
