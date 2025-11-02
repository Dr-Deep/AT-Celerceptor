package atfram

// Server Version wahrscheinlich:
// ForgeRock OpenAM (XUI) + OAuth2 (AM 13–14 era, evtl. AM 5.x rebrand).
// und da ist auch eine fette proxy irwo, die blockt ab und zu pfade

const (
	ALDITALK_URL      = "https://login.alditalk-kundenbetreuung.de"
	ALDITALK_BASE     = "/signin/json/realms/root/realms"
	ALDITALK_REALM    = "alditalk"
	ALDITALK_BASE_URI = ALDITALK_URL + ALDITALK_BASE + ALDITALK_REALM

	// AM Resource / Endpoints
	ALDITALK_AUTHENTICATE_EP = ALDITALK_BASE_URI + "/authenticate"
	ALDITALK_SESSIONS_EP     = ALDITALK_BASE_URI + "/sessions"
	ALDITALK_USERS_EP        = ALDITALK_BASE_URI + "/users"
	/*[
		authenticate,
		users,
		groups,
		agents,
		realms,
		dashboard,
		sessions,
		serverinfo/,
		users/{%s},
		applications,
		resourcetypes,
		policies,
		applicationtypes,
		conditiontypes,
		subjecttypes,
		subjectattributes,
		decisioncombiners,
		subjectattributes,
		client,
	]*/
)

/*
stateful client wird das hier
* funcs:
- Login()
- RefreshSession() ?
- Logout()
- GetUserInfo() ?
- GetDataVolume()
- datenvolumen nachbuchen
*/

func (c *Client) Login() error {
	// jze callbacks ausfüllen
	// loop:ausfüllen, abschicken => bis wir Response mit TokenID haben

	/*
		first req
		- fill callbacks
		- send callbacks
		- more callbacks? => fillcallbacks
		- oder TokenID?
	*/

	// first Request
	resp, err := c.authenticate()
	if err != nil {
		return err
	}

	for {
		// müssen wir weiter callbacks ausfüllen?
		if resp.TokenID != "" {
			//! logged in
		}

		// parsen
		requirements, err := c.getRequirements(resp.Callbacks)

		// ausfüllen
		err := c.submitRequirements(requirements)

		// senden
	}

	c.submitRequirements()

}

func (c *Client) getRequirements(callbacks []CallbackRaw) ([]Callback, error) {
	// parse last

	for idx, cb := range callbacks {

	}
}

func (c *Client) submitRequirements(callbacks []Callback) error {
	// geparste resp solven && senden
	// wenn kein tokenID: c.submitRequirements(newest_resp_callbacks)

}

func (c *Client) hasMoreRequirements() bool {
	// stateful status?
}

// https://openam.example.com:8443/openam/json/sessions/?_action=logout&tokenId=IRWAS
func (c *Client) Logout() {}

// validation: /json/sessions/AQIC5...?_action=validate

// /json/sessions/?_action=getTimeLeft&tokenId=IRWAS
// /json/sessions/?_action=getMaxIdle
// json/sessions/?_action=isActive&refresh=true&tokenId=IRWAS
