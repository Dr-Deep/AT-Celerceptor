package atfram

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	curdoc []Callback

	// authID?
	// clientID?
	// customerID?
	// etc

	c      *http.Client
	logger *logging.Logger
	atconf *AldiTalkConfig
}

func (c *Client) getRequirements(callbacks []CallbackRaw) ([]Callback, error) {
	var requirements []Callback
	for _, cbraw := range callbacks {
		cb, err := matchCallback(cbraw)
		if err != nil {
			return nil, err
		}

		requirements = append(requirements, cb)
	}

	return requirements, nil
}

//c.logger.Debug("SEND: ", req.URL.String(), fmt.Sprintf("%#v\n\n", req))
//c.logger.Debug("RECEIVED:", fmt.Sprintf("%#v\n", rawresp), "BODY:", readToBuf(rawresp.Body).String(), "\n\n")

// geparste resp solven && senden
// wenn kein tokenID: c.submitRequirements(newest_resp_callbacks)

/*
? mit json body bzw filled callbacks
* /authenticate
"authId": "",
"callbacks": [],
"stage": "",
"header": ""
*/
func (c *Client) submitRequirements(callbacks []Callback) (*Response, error) {

	// marshal
	jsonBody, err := json.Marshal(callbacks)
	if err != nil {
		return nil, err
	}

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

func (c *Client) solveRequirements(callbacks []Callback) error {
	for _, cb := range callbacks {
		c.logger.Info("solving", cb.Type().String())

		switch cb.Type() {
		/*case CHOICE_CALLBACK:*/

		case CONFIRMATION_CALLBACK:
			c.solveConfirmationCallback(cb)

		case HIDDEN_VALUE_CALLBACK:
			c.solveHiddenValueCallback(cb)

		/*case HTTP_CALLBACK:*/

		/*case LANGUAGE_CALLBACK:*/

		case NAME_CALLBACK:
			c.solveNameCallback(cb)

		case PASSWORD_CALLBACK:
			c.solvePasswordCallback(cb)

		/*case REDIRECT_CALLBACK:*/

		/*case SCRIPT_TEXT_OUTPUT_CALLBACK:*/

		/*case TEXT_INPUT_CALLBACK:*/

		case TEXT_OUTPUT_CALLBACK:

		/*case X509_CERT_CALLBACK:*/

		default:
			err := fmt.Errorf("cant solve requirements: %w", ErrAldiTalkClientUnknownCallback)
			c.logger.Error(err.Error(), fmt.Sprintf("%#v", cb))
			return err
		}
	}

	return nil
}
