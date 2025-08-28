package alditalk

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/139.0"
	BaseURL   = "https://login.alditalk-kundenbetreuung.de"

	DashboardURI = "/portal/auth/uebersicht/"

	LoginURI = "/signin/XUI/#login/"

	RealmURI    = "/signin/json/realms/root/realms/alditalk"
	AuthURI     = BaseURL + RealmURI + "/authenticate"
	SessionsURI = BaseURL + "/signin/json/sessions"
	UsersURI    = BaseURL + RealmURI + "/users"
)

/*
1. Start – Login-Anfrage
POST → /signin/json/realms/root/realms/alditalk/authenticate
Hier wird der Login-Prozess gestartet.
Wahrscheinlich werden Benutzername (Rufnummer/E-Mail) und Passwort übertragen.
Antwort: 200 OK mit JSON (enthält Challenge/Session-Infos).


2. Session anlegen
POST → /signin/json/sessions
Server legt eine Session an.
Antwort: 200 OK, JSON mit Session-ID / Token.


3. Benutzer-Check
POST → /signin/json/realms/root/realms/alditalk/users
Abgleich, ob Benutzer in diesem Realm existiert.
Antwort: 200 OK, JSON mit User-Daten.


4. Userinfo laden
GET → /signin/json/realms/root/realms/alditalk/users/...
Details zum Account werden geholt.
Antwort: 200 OK, JSON mit User-Infos.


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

type AldiTalkSession struct {
	username   string
	password   string
	customerID string

	client *http.Client
}

func NewAldiTalkSession(username, password string) (*AldiTalkSession, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	var client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 30,
	}

	var s = &AldiTalkSession{
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
func (s *AldiTalkSession) login() error {}

/*
// Login Flow via ForgeRock
func (s *AldiSession) login() error {
	// 1. Initial Challenge
	req, _ := http.NewRequest("POST",
		"https://login.alditalk-kundenbetreuung.de/signin/json/authenticate",
		bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var challenge map[string]any
	if err := json.Unmarshal(body, &challenge); err != nil {
		return err
	}

	authId, ok := challenge["authId"].(string)
	if !ok {
		return fmt.Errorf("kein authId in Challenge gefunden")
	}

	// 2. Credentials zurücksenden
	payload := map[string]any{
		"authId": authId,
		"callbacks": []any{
			map[string]any{
				"type": "NameCallback",
				"input": []map[string]string{
					{"name": "IDToken1", "value": s.User},
				},
			},
			map[string]any{
				"type": "PasswordCallback",
				"input": []map[string]string{
					{"name": "IDToken2", "value": s.Pass},
				},
			},
		},
	}

	data, _ := json.Marshal(payload)

	req2, _ := http.NewRequest("POST",
		"https://login.alditalk-kundenbetreuung.de/signin/json/authenticate",
		bytes.NewBuffer(data))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := s.client.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	out, _ := io.ReadAll(resp2.Body)

	// Prüfen ob Login erfolgreich
	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return err
	}

	if _, ok := result["tokenId"]; !ok {
		return fmt.Errorf("login fehlgeschlagen: %s", string(out))
	}

	// 3. Jetzt /users abfragen → CID holen
	return s.fetchCID()
}

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
