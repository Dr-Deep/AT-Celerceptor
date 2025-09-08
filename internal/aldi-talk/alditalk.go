package alditalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	//DashboardURI = "/portal/auth/uebersicht/"
	//LoginURI     = "/signin/XUI/#login/"

	ALDITALK_BASE_URL              = "https://login.alditalk-kundenbetreuung.de"
	ALDITALK_ENDPOINT_REALM_URI    = "/signin/json/realms/root/realms/alditalk"
	ALDITALK_ENDPOINT_AUTH_URI     = ALDITALK_ENDPOINT_REALM_URI + "/authenticate"
	ALDITALK_ENDPOINT_SESSIONS_URI = "/signin/json/sessions"
	ALDITALK_ENDPOINT_USERS_URI    = ALDITALK_ENDPOINT_REALM_URI + "/users"
)

var (
	customerIDRegexp = regexp.MustCompile(`^C-\d{10}$`) // C-0123456789

	defaultHeaders = map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/139.0",
	}

	//? Errors
)

type AldiTalk_Auth_Challenge_Callbacks struct {
	Type   string                             `json:"type"`
	Input  []AldiTalk_Auth_Challenge_Callback `json:"input"`
	Output []AldiTalk_Auth_Challenge_Callback `json:"output"`
	Id     int                                `json:"_id"`
}

// bei ConfirmationCallback kann auch value:[]string
type AldiTalk_Auth_Challenge_Callback struct {
	Name  string           `json:"name"`
	Value *json.RawMessage `json:"value"`
}

func (c *AldiTalk_Auth_Challenge_Callback) GetValueAsString() (string, error) {
	if c.Value == nil {
		return "", fmt.Errorf("value is nil")
	}

	var valueString string
	if err := json.Unmarshal(*c.Value, &valueString); err != nil {
		return "", err
	}

	return valueString, nil
}

func (c *AldiTalk_Auth_Challenge_Callback) SetValueAsString(value string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	raw := json.RawMessage(data)
	c.Value = &raw

	return nil
}

func (c *AldiTalk_Auth_Challenge_Callback) GetValueAsArray() ([]string, error) {
	if c.Value == nil {
		return nil, fmt.Errorf("value is nil")
	}

	var valueArray []string
	if err := json.Unmarshal(*c.Value, &valueArray); err != nil {
		return nil, err
	}

	return valueArray, nil
}

func (c *AldiTalk_Auth_Challenge_Callback) SetValueAsArray(value []string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	raw := json.RawMessage(data)
	c.Value = &raw

	return nil
}

type AldiTalk_Auth_Challenge struct {
	AuthID    string                              `json:"authId"`
	Callbacks []AldiTalk_Auth_Challenge_Callbacks `json:"callbacks"`
	Stage     string                              `json:"stage"`
	Header    string                              `json:"header"`
}

/*

- ate
* config?
* login session '/login'
	* cookie banner
	* cookies speichern
	* einloggen
* datenvolumen? (int+einheit)
* gotify (success,fail)
* das alles im intervall
*/

func setDefaultHeadersForReq(req *http.Request) {
	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}
}

func readRespToBuf(r io.ReadCloser) (*bytes.Buffer, error) {
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	//? io.LimitReader()

	return bytes.NewBuffer(out), r.Close()
}
