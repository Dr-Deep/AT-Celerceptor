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

func (cr CallbackRaw) Marshal() {
	//json.Marshal(cr)
}

type Request struct {
	AuthID    string     `json:"authId"` // JSON Web Token (JWT)
	Callbacks []Callback `json:"callbacks"`
	//? stage etc
	//? header
}

type Response struct {
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
