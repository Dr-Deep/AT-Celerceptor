package atfram

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
func (c *Client) solveConfirmationCallback(callbacks []Callback, _cb Callback) error {
	//? eif so lassen

	return nil
}

func (c *Client) solveHiddenValueCallback(callbacks []Callback, _cb Callback) error {
	_ = _cb.(*HiddenValueCallback)

	//?

	return nil

}

func (c *Client) solveNameCallback(callbacks []Callback, _cb Callback) error {
	_cb.(*NameCallback).SetUsername(c.atconf.Username)

	return nil
}

func (c *Client) solvePasswordCallback(callbacks []Callback, _cb Callback) error {
	_cb.(*PasswordCallback).SetPassword(c.atconf.Password)

	return nil
}

func (c *Client) solveTextOutputCallback(callbacks []Callback, _cb Callback) error {
	cb := _cb.(*TextOutputCallback)

	// Proof of Work
	var proofOfWorkHash string
	{
		powJs := cb.GetMessage("message") // => pow script
		if powJs == "" {
			return ErrAldiTalkCallbackEmptyPoWScript
		}

		var (
			proofOfWorkUUID       = aldiTalk_PoW_UUID_RE.FindString(powJs)
			proofOfWorkDifficulty = aldiTalk_PoW_Difficulty_RE.FindString(powJs)
		)

		hash, err := GetProofOfWorkHash(proofOfWorkUUID, proofOfWorkDifficulty)
		if err != nil {
			return err
		}

		proofOfWorkHash = hash
	}

	// Submit into HiddenValueCallback
	{
		for _, __cb := range callbacks {
			if __cb.Type() != HIDDEN_VALUE_CALLBACK {
				continue
			}

			hvcb := _cb.(*HiddenValueCallback)

			//
			if hvcb.GetValue() == "proofOfWorkNonce" {
				hvcb.SetValue(
					proofOfWorkHash,
				)

				break
			}
		}
	}

	return nil
}
