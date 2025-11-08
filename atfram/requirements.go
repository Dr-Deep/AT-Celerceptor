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
func (c *Client) solveConfirmationCallback(cb Callback) {
	//? eif so lassen
}

func (c *Client) solveHiddenValueCallback(cb Callback) {
	hvcb, ok := cb.(*HiddenValueCallback)
	if !ok {
		return
	}

	//! POW
	if hvcb.GetValue() == "proofOfWorkNonce" {
		hvcb.SetValue("POW HASH")
	}

	//else irwas

}

func (c *Client) solveNameCallback(cb Callback) {
	ncb, ok := cb.(*NameCallback)
	if ok {
		ncb.SetUsername(c.atconf.Username)
	}
}

func (c *Client) solvePasswordCallback(cb Callback) {
	pcb, ok := cb.(*PasswordCallback)
	if ok {
		pcb.SetPassword(c.atconf.Password)
	}
}

func (c *Client) solveTextOutputCallback(cb Callback) {
	textOutcb, ok := cb.(*TextOutputCallback)
	if !ok {
		return
	}

	//
	powJs := textOutcb.GetMessage("message") // => pow script
	if powJs == "" {
		return
	}

	var (
		uuid       = aldiTalk_PoW_UUID_RE.FindString(powJs)
		difficulty = aldiTalk_PoW_Difficulty_RE.FindString(powJs)
	)

	nonce := proofOfWork(uuid, difficulty)

	/*
		in HiddenValueCallback rein
	*/

	// alle hiddenvalue callbacks durchgehen bis wir eins haben mit
	// cb.GetValue() == "proofOfWorkNonce"
	// dann cb.SetValue() auf nonce
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
