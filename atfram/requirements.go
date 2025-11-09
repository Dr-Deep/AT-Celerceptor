package atfram

import "fmt"

// idk was das tut
func (c *Client) solveConfirmationCallback(callbacks []Callback, _cb Callback) error {
	cb := _cb.(*ConfirmationCallback)
	cb.Inputs[0].Value = "2"

	return nil
}

// idk was das tut aber so passierts
func (c *Client) solveHiddenValueCallback(callbacks []Callback, _cb Callback) error {
	cb := _cb.(*HiddenValueCallback)

	// für alle außer pow
	if cb.GetValue() != "proofOfWorkNonce" {
		cb.Inputs[0].Value = cb.Outputs[1].Value
	}

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
	var (
		proofOfWorkNonce      string
		proofOfWorkUUID       string
		proofOfWorkDifficulty string
	)
	{
		powJs := cb.GetMessage("message") // => pow script
		if powJs == "" {
			return ErrAldiTalkCallbackEmptyPoWScript
		}

		// Find Vars
		var (
			matchUUID = aldiTalk_PoW_UUID_RE.FindStringSubmatch(powJs)
			matchDiff = aldiTalk_PoW_Difficulty_RE.FindStringSubmatch(powJs)
		)
		if matchUUID == nil || matchDiff == nil {
			return ErrAldiTalkCallbackPoWNoMatch
		}

		//
		proofOfWorkUUID = matchUUID[1]
		proofOfWorkDifficulty = matchDiff[1]
		//proofOfWorkDifficulty = cb.GetMessage("messageType")

		nonce, err := GetProofOfWorkNonce(proofOfWorkUUID, proofOfWorkDifficulty)
		if err != nil {
			return err
		}

		proofOfWorkNonce = nonce
	}

	// Submit into HiddenValueCallback
	{
		for _, __cb := range callbacks {
			if __cb.Type() != HIDDEN_VALUE_CALLBACK {
				continue
			}

			hvcb := __cb.(*HiddenValueCallback)
			//
			if hvcb.GetValue() == "proofOfWorkNonce" {
				hvcb.SetValue(
					proofOfWorkNonce,
				)

				c.logger.Info(
					"Proof of Work gelöst und gesetzt",
					fmt.Sprintf("diff: %s", proofOfWorkDifficulty),
					proofOfWorkNonce,
				)

				c.logger.Debug(hvcb.Type().String(), fmt.Sprintf("%#v", hvcb))

				break
			}
		}
	}

	return nil
}
