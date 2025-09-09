package alditalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ALDITALK_CALLBACK_Name         = "NameCallback"        // Benutzername / Rufnummer
	ALDITALK_CALLBACK_Password     = "PasswordCallback"    // Passwort
	ALDITALK_CALLBACK_TextOutput   = "TextOutputCallback"  // Nur Text für den Benutzer
	ALDITALK_CALLBACK_HiddenValue  = "HiddenValueCallback" // Proof-of-Work, Nonces, technische Parameter
	ALDITALK_CALLBACK_Confirmation = "ConfirmationCallback"
	//ALDITALK_CALLBACK_Choice       = "ChoiceCallback" // Auswahloptionen
)

/*
# Allgemeiner Ablauf bei ForgeRock-AM Login

* Initialer POST an /json/realms/root/authenticate?...
→ Server schickt dir die erste Challenge mit callbacks.

* Client wertet callbacks aus, füllt die Inputs (IDTokenX) mit passenden Werten (z. B. Benutzername, Passwort, hier Proof-of-Work Nonce).

* Antwort zurück an AM im gleichen JSON-Format.

* Wiederholung, bis entweder:
- Du ein tokenId oder successUrl bekommst (Login erfolgreich, jetzt hast du die Session = CID), oder
- der Flow fehlschlägt (Fehler/Abbruch).
*/

/* Login Flow via ForgeRock
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
*/

/*
- Solve ForgeRock Challenges -> bis wir:

	{
	"tokenId":"1FaaBylL0qRbIFvBsfI3bnAGrjw.*AAJTSQACMDMAAlNLABw5Zk9FS2hmaFVmcGFaUENtcG15a1o2dCtMSWs9AAR0eXBlAANDVFMAAlMxAAIwNA..*",
	"successUrl":"https://www.alditalk-kundenbetreuung.de/",
	"realm":"/alditalk",
	}
*/
func (s *AldiTalkSession) authenticate() (any, error) {
	// get challenges
	// solve challenges
	// send challenges (?token oder nochmal solven)

	var (
		currentChallenges *AldiTalk_Auth_Challenge
		err               error
	)
	for {
		// Request/Respond
		currentChallenges, err = s.authGetChallenges(currentChallenges)
		if err != nil {
			return nil, err
		}

		s.Logger.Debug("Auth Stage", currentChallenges.Stage)

		//? if success bzw token ist da
		if true == false {
			break
		}

		// Solve
		if err := s.authSolveChallenges(currentChallenges); err != nil {
			return nil, err
		}
	}
}

// Get ForgeRock Login Challenges
func (s *AldiTalkSession) authGetChallenges(challenges *AldiTalk_Auth_Challenge) (*AldiTalk_Auth_Challenge, error) {
	var jsonBody []byte
	if challenges == nil {
		jsonBody = []byte(`{}`)
	} else {
		// marshal
		_jsonBody, err := json.MarshalIndent(challenges, "", "  ")
		if err != nil {
			return nil, err
		}
		jsonBody = _jsonBody
	}

	// HTTP POST
	req, err := http.NewRequest(
		"POST",
		ALDITALK_BASE_URL+ALDITALK_ENDPOINT_AUTH_URI,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	// Header
	setDefaultHeadersForReq(req)

	// Do Request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	// valid Response?
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status-code is not 200")
	}

	respBuf, err := readToBuf(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResp AldiTalk_Auth_Challenge
	if err := json.NewDecoder(respBuf).Decode(&authResp); err != nil {
		s.Logger.Debug("Request", string(jsonBody))
		s.Logger.Debug("Response", respBuf.String())
		return nil, err
	}

	if authResp.AuthID == "" {
		return nil, fmt.Errorf("keine authId in response")
	}

	return &authResp, nil
}

// Solve ForgeRock Login Challenges
func (s *AldiTalkSession) authSolveChallenges(challenges *AldiTalk_Auth_Challenge) error {

	for i, challenge := range challenges.Callbacks {
		var idx = fmt.Sprintf("(%v/%v): %v", i+1, len(challenges.Callbacks), challenge.Id)

		switch challenge.Type {
		case ALDITALK_CALLBACK_Name:
			s.Logger.Debug(idx, "Got", "NameCallback", fmt.Sprintf("%#v", challenge))
			if err := s.authSolveChallengeName(&challenge); err != nil {
				s.Logger.Error("ALDITALK_CALLBACK_Name", err.Error())
				return err
			}

		case ALDITALK_CALLBACK_Password:
			s.Logger.Debug(idx, "Got", "PasswordCallback", fmt.Sprintf("%v", challenge))
			if err := s.authSolveChallengePassword(&challenge); err != nil {
				s.Logger.Error("ALDITALK_CALLBACK_Password", err.Error())
				return err
			}

		case ALDITALK_CALLBACK_TextOutput:
			s.Logger.Debug(idx, "Got", "TextOutputCallback", fmt.Sprintf("%v", challenge))
			if err := s.authSolveChallengeTextOutput(&challenge); err != nil {
				s.Logger.Error("ALDITALK_CALLBACK_TextOutput", err.Error())
				return err
			}

		case ALDITALK_CALLBACK_HiddenValue:
			s.Logger.Debug(idx, "Got", "HiddenValueCallback", fmt.Sprintf("%v", challenge))
			if err := s.authSolveChallengeHiddenValue(&challenge); err != nil {
				s.Logger.Error("ALDITALK_CALLBACK_HiddenValue", err.Error())
				return err
			}

		case ALDITALK_CALLBACK_Confirmation:
			//! bei ConfirmationCallback kann auch value:[]string
			s.Logger.Debug(idx, "Got", "ConfirmationCallback", fmt.Sprintf("%v", challenge))
			//s.authSolveChallengeConfirmation()

		default:
			return fmt.Errorf(idx, "unknown aldi-talk(forgerock-am) challenge callback")
		}
	}

	return nil
}

/*
"type": "NameCallback",
"output": [

	{
	    "name": "prompt",
	    "value": "custom.alditalk.loginuserbasic.login|custom.alditalk.loginuserbasic.loginplaceholder"
	}

],
"input": [

	{
	    "name": "IDToken3",
	    "value": "TELNUM/Username"
	}

],
"_id": 2
*/
func (s *AldiTalkSession) authSolveChallengeName(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
	return fillCallback(
		challenge,
		// Input
		map[string]string{
			"IDToken3": s.username,
		},
		// Output
		map[string]string{
			"prompt": "custom.alditalk.loginuserbasic.login|custom.alditalk.loginuserbasic.loginplaceholder",
		},
	)
}

/*
"type": "PasswordCallback",
"output": [

	{
	    "name": "prompt",
	    "value": "custom.alditalk.loginuserbasic.password|custom.alditalk.loginuserbasic.passwordplaceholder"
	}

],
"input": [

	{
	    "name": "IDToken4",
	    "value": "PASSWORD"
	}

],
"_id": 3
*/
func (s *AldiTalkSession) authSolveChallengePassword(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
	return fillCallback(
		challenge,
		// Input
		map[string]string{
			"IDToken4": s.password,
		},
		// Output
		map[string]string{
			"prompt": "custom.alditalk.loginuserbasic.password|custom.alditalk.loginuserbasic.passwordplaceholder",
		},
	)
}

func (s *AldiTalkSession) authSolveChallengeTextOutput(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
	/*
			 "type": "TextOutputCallback",
		            "output": [
		                {
		                    "name": "message",
		                    "value": "function startProofOfWork(uuid, difficulty, onProofOfWorkSuccess)"
		                },
		                {
		                    "name": "messageType",
		                    "value": "4"
		                }
		            ],
		            "_id": 1
	*/
	return nil
}

func (s *AldiTalkSession) authSolveChallengeHiddenValue(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
	/*
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
	*/
	return nil
}

func (s *AldiTalkSession) authSolveChallengeConfirmation(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
	/*
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
	*/
	return nil
}
