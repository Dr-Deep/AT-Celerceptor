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

// ? 	url.URL
type Client struct {
	authID  string
	tokenID string

	// clientID?
	// customerID?
	// etc

	c      *http.Client
	logger *logging.Logger
	atconf *AldiTalkConfig
}

// SuccessURL:"https://www.alditalk-kundenbetreuung.de/",

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

func (c *Client) submitRequirements(callbacks []Callback) (*Response, error) {
	var jsonBody = []byte{'{', '}'}

	if len(callbacks) > 0 {
		// marshal
		_jsonBody, err := json.Marshal(
			&Request{
				AuthID:    c.authID,
				Callbacks: callbacks,
			})
		if err != nil {
			return nil, err
		}
		jsonBody = _jsonBody
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
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)
	}

	var resp Response
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) solveRequirements(callbacks []Callback) error {
	for idx, cb := range callbacks {

		var err error
		switch cb.Type() {
		/*case CHOICE_CALLBACK:*/

		case CONFIRMATION_CALLBACK:
			err = c.solveConfirmationCallback(callbacks, cb)

		case HIDDEN_VALUE_CALLBACK:
			err = c.solveHiddenValueCallback(callbacks, cb) //?

		/*case HTTP_CALLBACK:*/

		/*case LANGUAGE_CALLBACK:*/

		case NAME_CALLBACK:
			err = c.solveNameCallback(callbacks, cb)

		case PASSWORD_CALLBACK:
			err = c.solvePasswordCallback(callbacks, cb)

		/*case REDIRECT_CALLBACK:*/

		/*case SCRIPT_TEXT_OUTPUT_CALLBACK:*/

		/*case TEXT_INPUT_CALLBACK:*/

		case TEXT_OUTPUT_CALLBACK:
			err = c.solveTextOutputCallback(callbacks, cb)

		/*case X509_CERT_CALLBACK:*/

		default:
			err := fmt.Errorf("cant solve requirements: %w", ErrAldiTalkClientUnknownCallback)
			c.logger.Error(err.Error(), fmt.Sprintf("%#v", cb))
			return err
		}

		if err != nil {
			return err
		}

		c.logger.Info("solved", fmt.Sprintf("(%v/%v)", idx+1, len(callbacks)), cb.Type().String())
	}

	return nil
}
