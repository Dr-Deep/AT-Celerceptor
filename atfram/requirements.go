package atfram

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *Client) submitRequirements(callbacks []Callback) (*FramResponse, error) {
	var jsonBody []byte

	if len(callbacks) > 0 {
		// marshal
		_jsonBody, err := json.Marshal(
			&FramRequest{
				AuthID:    c.authID,
				Callbacks: callbacks,
			})
		if err != nil {
			return nil, err
		}
		jsonBody = _jsonBody
	} else {
		jsonBody = []byte{'{', '}'}
	}

	// HTTP POST
	rawresp, err := c.DoHttpRequest(
		http.MethodPost,
		ALDITALK_AUTHENTICATE_EP,
		jsonBody,
	)
	if err != nil {
		return nil, err
	}

	// is resp valid
	switch rawresp.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)

	case http.StatusOK:
		break

	default:
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)
	}

	// decode resp
	var resp FramResponse
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		c.logger.Debug(_readToBuf(rawresp.Body).String())
		return nil, err
	}

	return &resp, nil
}

func (c *Client) solveRequirements(callbacks []Callback) error {
	var matchCallback = func(cb Callback) (err error) {

		switch cb.Type() {
		/*case CHOICE_CALLBACK:*/

		case CONFIRMATION_CALLBACK:
			err = c.solveConfirmationCallback(callbacks, cb)

		case HIDDEN_VALUE_CALLBACK:
			err = c.solveHiddenValueCallback(callbacks, cb)

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
			return fmt.Errorf(
				"cant solve requirements: %w",
				ErrAldiTalkClientUnknownCallback,
			)
		}

		return err
	}

	// * solve

	for idx, cb := range callbacks {
		if err := matchCallback(cb); err != nil {
			return err
		}

		c.logger.Info(
			fmt.Sprintf("(%v/%v)", idx+1, len(callbacks)),
			"solved requirement",
			cb.Type().String(),
		)
	}

	return nil
}
