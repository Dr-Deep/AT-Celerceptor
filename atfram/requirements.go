package atfram

import "fmt"

func (c *Client) test(requirements []Callback) {
	for _, v := range requirements {
		c.logger.Info(fmt.Sprintf("%#v", v))
	}

	c.logger.Fatal("test() !")
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

func (c *Client) solveRequirements(callbacks []Callback) error {

	for _, cb := range callbacks {

		switch cb.Type() {
		//
		default:
			//! err
		}
	}
}

func (c *Client) submitRequirements(callbacks []Callback) error {
	// geparste resp solven && senden
	// wenn kein tokenID: c.submitRequirements(newest_resp_callbacks)

	return nil
}

func (c *Client) hasMoreRequirements() bool {
	// stateful status?
	return false
}
