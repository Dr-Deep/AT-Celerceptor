package atfram

import "encoding/json"

const (
	UNKNOWN_CALLBACK CallbackType = iota
	CHOICE_CALLBACK
	CONFIRMATION_CALLBACK
	HIDDEN_VALUE_CALLBACK
	HTTP_CALLBACK
	LANGUAGE_CALLBACK
	NAME_CALLBACK
	PASSWORD_CALLBACK
	REDIRECT_CALLBACK
	SCRIPT_TEXT_OUTPUT_CALLBACK
	TEXT_INPUT_CALLBACK
	TEXT_OUTPUT_CALLBACK
	X509_CERT_CALLBACK
)

type CallbackType int

func (ct CallbackType) String() string {
	for k, v := range callbackRegistry {
		if v == ct {
			return k
		}
	}

	return "UNKNOWN"
}

type Callback interface {
	Type() CallbackType
	Prompt() string
}

type CallbackInput struct {
	Name  string     `json:"name"`
	Value json.Token `json:"value"`
}

type CallbackOutput struct {
	Name  string     `json:"name"`
	Value json.Token `json:"value"`
}

type CallbackRaw struct {
	ID           int              `json:"_id"`
	CallbackType string           `json:"type"`
	Inputs       []CallbackInput  `json:"input"`
	Outputs      []CallbackOutput `json:"output"`
}

type FramRequest struct {
	AuthID    string     `json:"authId"` // JSON Web Token (JWT)
	Callbacks []Callback `json:"callbacks"`
}

/*
! Failure:
{
	"code":401,
	"reason":"Unauthorized",
	"message":"Invalid Password!!"
}
* Success:
{
	"tokenId": "AQIC5wM2...U3MTE4NA..*",
	"successUrl": "/openam/console"
}
*/
type FramResponse struct {
	AuthID    string        `json:"authId"` // JSON Web Token (JWT)
	Template  string        `json:"template"`
	Stage     string        `json:"stage"`
	Callbacks []CallbackRaw `json:"callbacks"`

	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	TokenID string `json:"tokenId"`

	SuccessURL string `json:"successUrl"`
	FailureURL string `json:"failureUrl"`
}

type FramSessionInfo struct {
	UserID                   string `json:"username"`
	UniversalID              string `json:"universalId"`
	LatestAccessTime         string `json:"latestAccessTime"`
	MaxIdleExpirationTime    string `json:"maxIdleExpirationTime"`
	MaxSessionExpirationTime string `json:"maxSessionExpirationTime"`

	Properties struct {
		AMCtxID   string `json:"AMCtxId"`
		AuthLevel string `json:"auth_lvl"`
		Acr       string `json:"acr"`
	}
}

type FramUserInfo struct {
	UserID          string   `json:"username"`
	UID             []string `json:"uid"`
	TelephoneNumber []string `json:"telephoneNumber"`
	GivenName       []string `json:"givenName"`
	Surname         []string `json:"sn"`
	Roles           []string `json:"roles"`
}

/*
{
	"id":"A-IIIIIII",
	"realm":"/alditalk",
	"dn":"id=A-IIIIIII,ou=user,o=alditalk,ou=services,ou=am-config",
	"successURL":"https://www.alditalk-kundenbetreuung.de/",
	"fullLoginURL":"/signin/UI/Login?realm=%2Falditalk"
}
*/
type FramUIDFromSession struct {
	UserID       string `json:"id"`
	DN           string `json:"dn"`
	SuccessURL   string `json:"successURL"`
	FullLoginURL string `json:"fullLoginURL"`
}

type UserInfo struct {
	SessionID  string // session_id für scs
	UserID     string // user_id=A-xxxxx
	CustomerID string // tef_customer_id=C-xxxxx
	IsLoggedIn bool   // is_logged=true
}

// ciamsessionid
