package atfram

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
 * Access Management
- ForgeRock OpenAM (XUI) + OAuth2 (AM 13–14 era, evtl. AM 5.x rebrand).
- und da ist auch eine fette proxy irwo, die blockt ab und zu pfade

* Selfcare Dashboard (scs/bff/* API's) => (Backend For Frontend)
- irein O2/Telefónica backend über das ich keine ahnung habe
- kümmert sich um: Kontostand / Minuten / Daten, Tarife / Optionen, Offers API, Account Overview, Customer API
*/

/*
	* Domains
	- "login.alditalk-kundenbetreuung.de"	<= Login
	- "www.alditalk-kundenbetreuung.de"		<= SuccessURL
	- "www.alditalk-kundenportal.de"
*/

const (
	ALDITALK_URL      = "https://login.alditalk-kundenbetreuung.de"
	ALDITALK_BASE     = "/signin/json/realms/root/realms"
	ALDITALK_REALM    = "/alditalk"
	ALDITALK_BASE_URI = ALDITALK_URL + ALDITALK_BASE + ALDITALK_REALM

	// AM Resource / Endpoints
	ALDITALK_AUTHENTICATE_EP      = ALDITALK_BASE_URI + "/authenticate"
	ALDITALK_USERS_EP             = ALDITALK_BASE_URI + "/users"
	ALDITALK_GROUPS_EP            = ALDITALK_BASE_URI + "/groups"
	ALDITALK_AGENTS_EP            = ALDITALK_BASE_URI + "/agents"
	ALDITALK_REALMS_EP            = ALDITALK_BASE_URI + "/realms"
	ALDITALK_DASHBOARD_EP         = ALDITALK_BASE_URI + "/dashboard"
	ALDITALK_SESSIONS_EP          = ALDITALK_BASE_URI + "/sessions"
	ALDITALK_SERVERINFO_EP        = ALDITALK_BASE_URI + "/serverinfo/*"
	ALDITALK_APPLICATIONS_EP      = ALDITALK_BASE_URI + "/applications"
	ALDITALK_RESOURCETYPES_EP     = ALDITALK_BASE_URI + "/resourcetypes"
	ALDITALK_POLICIES_EP          = ALDITALK_BASE_URI + "/policies"
	ALDITALK_APPLICATIONTYPES_EP  = ALDITALK_BASE_URI + "/applicationtypes"
	ALDITALK_CONDITIONTYPES_EP    = ALDITALK_BASE_URI + "/conditiontypes"
	ALDITALK_SUBJECTTYPES_EP      = ALDITALK_BASE_URI + "/subjecttypes"
	ALDITALK_SUBJECTATTRIBUTES_EP = ALDITALK_BASE_URI + "/subjectattributes"
	ALDITALK_DECISIONCOMBINERS_EP = ALDITALK_BASE_URI + "/decisioncombiners"
	ALDITALK_CLIENT_EP            = ALDITALK_BASE_URI + "/client"
)

/*
validation: /json/sessions/AQIC5...?_action=validate
/json/sessions/?_action=getTimeLeft&tokenId=IRWAS
/json/sessions/?_action=getMaxIdle
/json/sessions/?_action=isActive&refresh=true&tokenId=IRWAS
*/

// /signin/json/sessions?_action=getSessionInfo
func (c *Client) GetFramSessionInfo() (*FramSessionInfo, error) {
	rawresp, err := c.DoHttpRequest(
		http.MethodPost,
		ALDITALK_SESSIONS_EP+"?_action=getSessionInfo",
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	// is resp valid
	if rawresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)
	}

	// decode resp
	var resp FramSessionInfo
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// /signin/json/realms/root/realms/alditalk/users?_action=idFromSession
func (c *Client) GetFramUIDFromSession() (*FramUIDFromSession, error) {
	rawresp, err := c.DoHttpRequest(
		http.MethodGet,
		ALDITALK_USERS_EP+"?_action=idFromSession",
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	// is resp valid?
	if rawresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)
	}

	// decode resp
	var resp FramUIDFromSession
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, err
}

// /signin/json/realms/root/realms/alditalk/users/A-IIIIIII
func (c *Client) GetFramUserInfo(userID string) (*FramUserInfo, error) {
	rawresp, err := c.DoHttpRequest(
		http.MethodGet,
		ALDITALK_USERS_EP+"/"+userID,
		[]byte{},
	)
	if err != nil {
		return nil, err
	}

	// is resp valid?
	if rawresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrAldiTalkClientInvalidStatusCode, rawresp.StatusCode)
	}

	// decode resp
	var resp FramUserInfo
	if err := json.NewDecoder(rawresp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, err
}

// from cookies
func (c *Client) GetUserInfo() (*UserInfo, error) {

	//c.c.Jar.Cookies()

	c.logger.Info(fmt.Sprintf("%#v\n", c.c.Jar))

	return nil, nil
}

/*
5. Weiterleitung ins Kundenportal
GET → /signin/oauth2/authorize?...
Klassische OAuth2-Autorisierung.
Server antwortet mit 302 Found → Redirect.

6. Redirects
GET → /openid/response?...
Server schickt 302 Found und leitet auf das Kundenportal (kundenportal.de) um.

7. Portal-Login bestätigen
GET → /user/auth/account-overview.html
Erste Seite nach erfolgreichem Login.
Antwort: 302 Found → weitere Weiterleitung.

GET → /logged-in-home-page/
Dashboard/Home-Seite im Kundenportal.
Antwort: 302 Found.

GET → /user/auth/account-overview.html (erneut)
Seite lädt erfolgreich.
Antwort: 301 Moved Permanently → endgültige URL.

GET → /user/auth/account-overview/
Endgültige Zielseite.
Antwort: 200 OK

Titel: „ALDI TALK | Übersicht“
*/

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
