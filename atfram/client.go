package atfram

import (
	"net/http"

	"github.com/Dr-Deep/logging-go"
)

type Client struct {
	authID string

	c      *http.Client
	logger *logging.Logger
	atconf *AldiTalkConfig
}

// returns tokenID & error
func (c *Client) Login() (*FramResponse, error) {
	var currentRequirements = []Callback{}

	//? timeout/deadline
	for {
		// submit solved requirements
		resp, err := c.submitRequirements(
			currentRequirements,
		)
		if err != nil {
			c.logger.Error("couldnt submit requirements", err.Error())
			return nil, err
		}

		// we made it
		if resp.TokenID != "" {
			return resp, nil
		}

		c.authID = resp.AuthID

		//? requirements ODER Erfolg?
		requirements, err := c.getRequirements(resp.Callbacks)
		if err != nil {
			c.logger.Error("couldnt get requirements", err.Error())
			return nil, err
		}

		// solve
		if err := c.solveRequirements(requirements); err != nil {
			c.logger.Error("couldnt solve requirements", err.Error())
			return nil, err
		}

		currentRequirements = requirements
	}
}

// https://openam.example.com:8443/openam/json/sessions/?_action=logout&tokenId=IRWAS
func (c *Client) Logout() {}
