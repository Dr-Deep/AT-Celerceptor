package atfram

var (
	callbackRegistry = map[string]CallbackType{
		"ChoiceCallback":           CHOICE_CALLBACK,
		"ConfirmationCallback":     CONFIRMATION_CALLBACK,
		"HiddenValueCallback":      HIDDEN_VALUE_CALLBACK,
		"HttpCallback":             HTTP_CALLBACK,
		"LanguageCallback":         LANGUAGE_CALLBACK,
		"NameCallback":             NAME_CALLBACK,
		"PasswordCallback":         PASSWORD_CALLBACK,
		"RedirectCallback":         REDIRECT_CALLBACK,
		"ScriptTextOutputCallback": SCRIPT_TEXT_OUTPUT_CALLBACK,
		"TextInputCallback":        TEXT_INPUT_CALLBACK,
		"TextOutputCallback":       TEXT_OUTPUT_CALLBACK,
		"X509CertificateCallback":  X509_CERT_CALLBACK,
	}

	callbackMatcherRegistry = map[CallbackType]func(CallbackRaw) (Callback, error){
		UNKNOWN_CALLBACK:            matchUnknownCallback,
		CHOICE_CALLBACK:             match_ChoiceCallback,
		CONFIRMATION_CALLBACK:       match_ConfirmationCallback,
		HIDDEN_VALUE_CALLBACK:       match_HiddenValueCallback,
		HTTP_CALLBACK:               match_HttpCallback,
		LANGUAGE_CALLBACK:           match_LanguageCallback,
		NAME_CALLBACK:               match_NameCallback,
		PASSWORD_CALLBACK:           match_PasswordCallback,
		REDIRECT_CALLBACK:           match_RedirectCallback,
		SCRIPT_TEXT_OUTPUT_CALLBACK: match_ScriptTextOutputCallback,
		TEXT_INPUT_CALLBACK:         match_TextInputCallback,
		TEXT_OUTPUT_CALLBACK:        match_TextOutputCallback,
		X509_CERT_CALLBACK:          match_X509CertCallback,
	}
)

/*
 * type matcher
 */

// !
func match_ChoiceCallback(c CallbackRaw) (Callback, error) {}

// !
func match_ConfirmationCallback(c CallbackRaw) (Callback, error) {}

/*
{
            "type": "HiddenValueCallback",
            "output": [
                {
                    "name": "value",
                    "value": ""
                },
                {
                    "name": "id",
                    "value": "proofOfWorkNonce"
                }
            ],
            "input": [
                {
                    "name": "IDToken1",
                    "value": "proofOfWorkNonce"
                }
            ],
            "_id": 0
},
{
            "type": "HiddenValueCallback",
            "output": [
                {
                    "name": "value",
                    "value": ""
                },
                {
                    "name": "id",
                    "value": "proofOfWorkNonce"
                }
            ],
            "input": [
                {
                    "name": "IDToken1",
                    "value": "892"
                }
            ],
            "_id": 0
        },
*/
// !
func match_HiddenValueCallback(c CallbackRaw) (Callback, error) {}

// !
func match_NameCallback(c CallbackRaw) (Callback, error) {
	return &NameCallback{CallbackRaw: c}, nil
}

// !
func match_PasswordCallback(c CallbackRaw) (Callback, error) {}

// !
func match_TextInputCallback(c CallbackRaw) (Callback, error) {}

// !
func match_TextOutputCallback(c CallbackRaw) (Callback, error) {}

func match_RedirectCallback(c CallbackRaw) (Callback, error)         { return matchUnknownCallback(c) }
func match_ScriptTextOutputCallback(c CallbackRaw) (Callback, error) { return matchUnknownCallback(c) }
func match_HttpCallback(c CallbackRaw) (Callback, error)             { return matchUnknownCallback(c) }
func match_LanguageCallback(c CallbackRaw) (Callback, error)         { return matchUnknownCallback(c) }
func match_X509CertCallback(c CallbackRaw) (Callback, error)         { return matchUnknownCallback(c) }

// ! log?; err
func matchUnknownCallback(c CallbackRaw) (Callback, error) {}

func matchCallback(cb CallbackRaw) (Callback, error) {
	var cbType, ok = callbackRegistry[cb.CallbackType]
	if !ok {
		cbType = UNKNOWN_CALLBACK
	}

	return callbackMatcherRegistry[cbType](cb)
}

/*
 * Used to display a list of choices and retrieve the selected choice.
 */

type ChoiceCallback struct {
	CallbackRaw
}

func (cb *ChoiceCallback) Type() CallbackType {
	return CHOICE_CALLBACK
}

func (cb *ChoiceCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 *  Used to ask for a confirmation such as Yes, No, or Cancel and retrieve the selection.
 */

/*
{
            "type": "ConfirmationCallback",
            "output": [
                {
                    "name": "prompt",
                    "value": ""
                },
                {
                    "name": "messageType",
                    "value": 0
                },
                {
                    "name": "options",
                    "value": [
                        "custom.alditalk.loginuserbasic.loginWithoutPassword",
                        "custom.alditalk.loginuserbasic.registerbtn",
                        "custom.alditalk.loginuserbasic.loginbtn",
                        "custom.alditalk.loginuserbasic.forgetP"
                    ]
                },
                {
                    "name": "optionType",
                    "value": -1
                },
                {
                    "name": "defaultOption",
                    "value": 1
                }
            ],
            "input": [
                {
                    "name": "IDToken5",
                    "value": 0
                }
            ],
            "_id": 4
        },
*/

type ConfirmationCallback struct {
	CallbackRaw
}

func (cb *ConfirmationCallback) Type() CallbackType {
	return CONFIRMATION_CALLBACK
}

func (cb *ConfirmationCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to return form values that are not visually rendered to the end user.
 */

// hier POW
type HiddenValueCallback struct {
	CallbackRaw
}

/*
        {
            "type": "HiddenValueCallback",
            "output": [
                {
                    "name": "value",
                    "value": ""
                },
                {
                    "name": "id",
                    "value": "proofOfWorkNonce"
                }
            ],
            "input": [
                {
                    "name": "IDToken1",
                    "value": "proofOfWorkNonce"
                }
            ],
            "_id": 0
},
{
            "type": "HiddenValueCallback",
            "output": [
                {
                    "name": "value",
                    "value": ""
                },
                {
                    "name": "id",
                    "value": "proofOfWorkNonce"
                }
            ],
            "input": [
                {
                    "name": "IDToken1",
                    "value": "892"
                }
            ],
            "_id": 0
},
*/

//! welches array? was zu tuen ist, ist komplizierter wegen verschiedener HiddenValue types bzw Aufgaben

func (cb *HiddenValueCallback) Type() CallbackType {
	return HIDDEN_VALUE_CALLBACK
}

func (cb *HiddenValueCallback) Prompt() string {
	return cb.Outputs[0].Value
}

//func (cb *HiddenValueCallback) SetHiddenValue(s string)

/*
 * Used for HTTP handshake negotiations.
 */

type HttpCallback struct {
	CallbackRaw
}

func (cb *HttpCallback) Type() CallbackType {
	return HTTP_CALLBACK
}

func (cb *HttpCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to retrieve the locale for localizing text presented to the end user.
 */

type LanguageCallback struct {
	CallbackRaw
}

func (cb *LanguageCallback) Type() CallbackType {
	return LANGUAGE_CALLBACK
}

func (cb *LanguageCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to retrieve a name string.
 */

type NameCallback struct {
	CallbackRaw
}

func (cb *NameCallback) Type() CallbackType {
	return NAME_CALLBACK
}

func (cb *NameCallback) Prompt() string {
	return cb.Outputs[0].Value
}

func (cb *NameCallback) GetName() string {
	return cb.Inputs[0].Name
}
func (cb *NameCallback) SetUsername(s string) {
	cb.Inputs[0].Value = s
}

/*
 * Used to retrieve a password value.
 */

type PasswordCallback struct {
	CallbackRaw
}

func (cb *PasswordCallback) Type() CallbackType {
	return PASSWORD_CALLBACK
}

func (cb *PasswordCallback) Prompt() string {
	return cb.Outputs[0].Value
}

func (cb *PasswordCallback) GetName() string {
	return cb.Inputs[0].Name
}

func (cb *PasswordCallback) SetPassword(s string) {
	cb.Inputs[0].Value = s
}

/*
 * Used to redirect the client user-agent.
 */

type RedirectCallback struct {
	CallbackRaw
}

func (cb *RedirectCallback) Type() CallbackType {
	return REDIRECT_CALLBACK
}

func (cb *RedirectCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to insert a script into the page presented to the end user. The script can, for example, collect data about the user’s environment.
 */

type ScriptTextOutputCallback struct {
	CallbackRaw
}

func (cb *ScriptTextOutputCallback) Type() CallbackType {
	return SCRIPT_TEXT_OUTPUT_CALLBACK
}

func (cb *ScriptTextOutputCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to retrieve text input from the end user.
 */

type TextInputCallback struct {
	CallbackRaw
}

func (cb *TextInputCallback) Type() CallbackType {
	return TEXT_INPUT_CALLBACK
}

func (cb *TextInputCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to display a message to the end user.
 */

type TextOutputCallback struct {
	CallbackRaw
}

func (cb *TextOutputCallback) Type() CallbackType {
	return TEXT_OUTPUT_CALLBACK
}

func (cb *TextOutputCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 * Used to retrieve the content of an x.509 certificate.
 */

type X509CertCallback struct {
	CallbackRaw
}

func (cb *X509CertCallback) Type() CallbackType {
	return X509_CERT_CALLBACK
}

func (cb *X509CertCallback) Prompt() string {
	return cb.Outputs[0].Value
}

/*
 *
 */
