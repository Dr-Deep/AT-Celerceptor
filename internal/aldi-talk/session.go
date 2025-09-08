package alditalk

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/Dr-Deep/logging-go"
)

type AldiTalkSession struct {
	username   string
	password   string
	customerID string

	Logger *logging.Logger
	client *http.Client
}

func NewAldiTalkSession(logger *logging.Logger, username, password string) (*AldiTalkSession, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	var client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 30,
	}

	var s = &AldiTalkSession{
		Logger:   logger,
		client:   client,
		username: username,
		password: password,
	}

	// perform login instantly
	if err := s.login(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *AldiTalkSession) RefreshSession() error {
	return s.login()
}

// ForgeRock Login Flow
func (s *AldiTalkSession) login() error {
	challenge, err := s.authGetChallenge()
	if err != nil {
		s.Logger.Error("authGetChallenge:", err.Error())
		return err
	}

	err = s.authSolveChallenges(challenge)
	if err != nil {
		s.Logger.Error("authSolveChallenges:", err.Error())
		return err
	}

	/*
		// Zurück in JSON
		newJSON, _ := json.MarshalIndent(callbacks, "", "  ")
		fmt.Println(string(newJSON))
	*/

	os.Exit(0)

	var authID = challenge.AuthID
	_ = authID

	//s.newSession()
	//s.checkUser()
	//s.getUserInfo()
	return nil
}

/*
  - 2. Session anlegen

POST → /signin/json/sessions
Server legt eine Session an.
Antwort: 200 OK, JSON mit Session-ID / Token.
*/
func (s *AldiTalkSession) newSession() {

	/*
		POST muss 200 sein; JSON
		/signin/json/sessions?_action=getSessionInfo
		Query: "_action=getSessionInfo"
	*/

}

/*
  - 3. Benutzer-Check

POST → /signin/json/realms/root/realms/alditalk/users
Abgleich, ob Benutzer in diesem Realm existiert.
Antwort: 200 OK, JSON mit User-Daten.
*/
func (s *AldiTalkSession) checkUser() {}

/*
  - 4. Userinfo laden

GET → /signin/json/realms/root/realms/alditalk/users/...
Details zum Account werden geholt.
Antwort: 200 OK, JSON mit User-Infos.
*/
func (s *AldiTalkSession) getUserInfo() {}

/*
// Holt die CustomerID (CID) aus dem Backend
func (s *AldiSession) fetchCID() error {
	req, _ := http.NewRequest("GET",
		"https://www.alditalk-kundenbetreuung.de/users", nil)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if cid, ok := data["customerId"].(string); ok {
		s.CID = cid
		return nil
	}

	return fmt.Errorf("keine CID gefunden, response: %s", string(body))
}
*/

//func (s *AldiTalkSession)

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
