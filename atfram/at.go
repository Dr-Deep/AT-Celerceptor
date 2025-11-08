package atfram

import "fmt"

// Server Version wahrscheinlich:
// ForgeRock OpenAM (XUI) + OAuth2 (AM 13–14 era, evtl. AM 5.x rebrand).
// und da ist auch eine fette proxy irwo, die blockt ab und zu pfade

const (
	ALDITALK_URL      = "https://login.alditalk-kundenbetreuung.de"
	ALDITALK_BASE     = "/signin/json/realms/root/realms"
	ALDITALK_REALM    = "/alditalk"
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
	var currentRequirements = []Callback{}

	for {
		// submit solved requirements
		resp, err := c.submitRequirements(
			currentRequirements,
		)
		if err != nil {
			return err
		}

		fmt.Printf("\n%#v\n\n", resp)
		fmt.Println("INPUT?")
		fmt.Scanln()

		//? requirements ODER Erfolg?
		requirements, err := c.getRequirements(resp.Callbacks)
		if err != nil {
			return err
		}

		//!
		currentRequirements = requirements

		// solve
		if err := c.solveRequirements(currentRequirements); err != nil {
			return err
		}

	}

	return nil
}

// https://openam.example.com:8443/openam/json/sessions/?_action=logout&tokenId=IRWAS
func (c *Client) Logout() {}

// validation: /json/sessions/AQIC5...?_action=validate

// /json/sessions/?_action=getTimeLeft&tokenId=IRWAS
// /json/sessions/?_action=getMaxIdle
// /json/sessions/?_action=isActive&refresh=true&tokenId=IRWAS
