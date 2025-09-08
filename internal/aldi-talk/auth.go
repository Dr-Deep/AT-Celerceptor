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
*/

// Get ForgeRock Login Challenge
func (s *AldiTalkSession) authGetChallenge() (*AldiTalk_Auth_Challenge, error) {
	// HTTP POST
	req, err := http.NewRequest(
		"POST",
		ALDITALK_BASE_URL+ALDITALK_ENDPOINT_AUTH_URI,
		bytes.NewBuffer([]byte(`{}`)),
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

	respBuf, err := readRespToBuf(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResp AldiTalk_Auth_Challenge
	if err := json.NewDecoder(respBuf).Decode(&authResp); err != nil {
		s.Logger.Debug("received", respBuf.String())
		return nil, err
	}

	if authResp.AuthID == "" {
		return nil, fmt.Errorf("keine authId in response")
	}

	return &authResp, nil
}

// Solve ForgeRock Challenges -> bis wir:
/*
{
	"tokenId":"1FaaBylL0qRbIFvBsfI3bnAGrjw.*AAJTSQACMDMAAlNLABw5Zk9FS2hmaFVmcGFaUENtcG15a1o2dCtMSWs9AAR0eXBlAANDVFMAAlMxAAIwNA..*",
	"successUrl":"https://www.alditalk-kundenbetreuung.de/",
	"realm":"/alditalk",
}
*/
func (s *AldiTalkSession) authSolveChallenges(challenges *AldiTalk_Auth_Challenge) error {

	for i, challenge := range challenges.Callbacks {
		var idx = fmt.Sprintf("(%v/%v): %v", i+1, len(challenges.Callbacks), challenge.Id)

		//
		switch challenge.Type {
		case ALDITALK_CALLBACK_Name:
			s.Logger.Debug(idx, "Got", "NameCallback", fmt.Sprintf("%#v", challenge))
			return s.authSolveChallengeName(&challenge)

		case ALDITALK_CALLBACK_Password:
			s.Logger.Debug(idx, "Got", "PasswordCallback", fmt.Sprintf("%v", challenge))
			s.authSolveChallengePassword(&challenge)

		case ALDITALK_CALLBACK_TextOutput:
			s.Logger.Debug(idx, "Got", "TextOutputCallback", fmt.Sprintf("%v", challenge))
			s.authSolveChallengeTextOutput(&challenge)

		case ALDITALK_CALLBACK_HiddenValue:
			s.Logger.Debug(idx, "Got", "HiddenValueCallback", fmt.Sprintf("%v", challenge))
			s.authSolveChallengeHiddenValue(&challenge)

		case ALDITALK_CALLBACK_Confirmation:
			//! bei ConfirmationCallback kann auch value:[]string
			s.Logger.Debug(idx, "Got", "ConfirmationCallback", fmt.Sprintf("%v", challenge))
			//s.authSolveChallengeConfirmation()

		default:
			return fmt.Errorf(idx, "unknown aldi-talk(forgerock-am) challenge callback")
		}
	}

	// send solved challenges?

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
	// Muster abstimmen?
	if len(challenge.Input) != 1 && len(challenge.Output) != 1 {
		//! err: falsches muster
	}

	// Output Muster
	ov, err := challenge.Output[0].GetValueAsString()
	if err != nil {
		//!err
	}

	if challenge.Output[0].Name != "prompt" || ov != "custom.alditalk.loginuserbasic.login|custom.alditalk.loginuserbasic.loginplaceholder" {
		//! err
	}

	// Input Muster
	if challenge.Input[0].Name != "IDToken3" {
		//! err
	}

	// Solve Challenge
	challenge.Input[0].SetValueAsString(s.username)

	return nil
}

func (s *AldiTalkSession) authSolveChallengePassword(challenge *AldiTalk_Auth_Challenge_Callbacks) error {
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
		                    "value": ""
		                }
		            ],
		            "_id": 3
	*/
	return nil
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

/*
	func fillCallbacks(fr frResponse, user, pass string) ([]map[string]any, error) {
		filled := make([]map[string]any, 0, len(cbArr))
		for _, cb := range cbArr {
			m := map[string]any{
				"type": cb.Type,
			}
			// Kopiere Output (falls Server es erwartet)
			if len(cb.Output) > 0 {
				m["output"] = cb.Output
			}

			// Fülle bekannte Callback-Typen
			switch strings.ToLower(cb.Type) {
			case "namecallback":
				m["input"] = []frInput{{Name: "IDToken1", Value: user}}
			case "passwordcallback":
				m["input"] = []frInput{{Name: "IDToken2", Value: pass}}
			default:
				// Unbekannt → leere Inputs übernehmen, falls vorhanden
				if len(cb.Input) > 0 {
					m["input"] = cb.Input
				}
			}
			filled = append(filled, m)
		}
		return filled, nil
}
*/

/*
	// Schritt 2: Challenge beantworten (Name/Password)
	filledChallenge, err := fillCallbacks(fr1, username, password)
	if err != nil {
		log.Fatalf("Callbacks nicht füllbar: %v", err)
	}
	payload2 := map[string]any{
		"authId":    fr1.AuthID,
		"callbacks": filledChallenge,
	}
	authResp2, err := postJSON(c, authenticateEP, payload2)
	if err != nil {
		log.Fatalf("Authenticate (Step 2) fehlgeschlagen: %v", err)
	}
*/
