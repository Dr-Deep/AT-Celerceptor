package atfram

import (
	"fmt"
	"strings"
)

// idk was das tut
func (c *Client) solveConfirmationCallback(_ []Callback, _cb Callback) error {
	cb := _cb.(*ConfirmationCallback)
	cb.Inputs[0].Value = "2"

	return nil
}

// idk was das tut aber so passierts
func (c *Client) solveHiddenValueCallback(_ []Callback, _cb Callback) error {
	cb := _cb.(*HiddenValueCallback)

	// für alle außer pow
	if cb.GetValue() != "proofOfWorkNonce" {
		cb.Inputs[0].Value = cb.Outputs[1].Value
	}

	return nil

}

func (c *Client) solveNameCallback(_ []Callback, _cb Callback) error {
	_cb.(*NameCallback).SetUsername(c.atconf.Username)

	return nil
}

func (c *Client) solvePasswordCallback(_ []Callback, _cb Callback) error {
	_cb.(*PasswordCallback).SetPassword(c.atconf.Password)

	return nil
}

func (c *Client) solveTextOutputCallback(callbacks []Callback, _cb Callback) error {
	var (
		cb  = _cb.(*TextOutputCallback)
		msg = cb.GetMessage("message")
	)

	if len(msg) == 0 {
		return nil
	}

	switch {
	case strings.HasPrefix(msg, "custom.alditalk.common.error$"):
		/* custom.alditalk.common.error$custom.alditalk.accountLock.accountLockMsg */
		var errmsg = msg
		s, oke := strings.CutPrefix(msg, "custom.alditalk.common.error$")
		if oke {
			errmsg = s
		}

		c.logger.Error("ALDITALK ERROR", errmsg)

		return fmt.Errorf(
			"%w: %s",
			ErrAldiTalkCustomErr,
			errmsg,
		)

	case strings.HasPrefix(msg, "function startProofOfWork"):
		if err := c.solvePoWCallback(callbacks, cb); err != nil {
			return err
		}

		return nil

	default:
		c.logger.Info("ALDITALK MESSAGE", msg)
		return nil
	}
}

func (c *Client) solvePoWCallback(callbacks []Callback, cb *TextOutputCallback) error {
	var (
		proofOfWorkNonce      string
		proofOfWorkUUID       string
		proofOfWorkDifficulty string
	)

	// * Solve ProofOfWork
	{
		powJsScript := cb.GetMessage("message") // => pow script

		// Find Vars
		var (
			matchUUID = aldiTalk_PoW_UUID_RE.FindStringSubmatch(powJsScript)
			matchDiff = aldiTalk_PoW_Difficulty_RE.FindStringSubmatch(powJsScript)
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

	// * Submit into HiddenValueCallback
	{
		for _, _cb := range callbacks {
			if _cb.Type() != HIDDEN_VALUE_CALLBACK {
				continue
			}

			hvcb := _cb.(*HiddenValueCallback)

			if hvcb.GetValue() == "proofOfWorkNonce" {
				hvcb.SetValue(
					proofOfWorkNonce,
				)

				c.logger.Info(
					"ProofOfWork solved",
					fmt.Sprintf("Difficulty: %s", proofOfWorkDifficulty),
					fmt.Sprintf("Nonce: %s", proofOfWorkNonce),
				)

				return nil
			}
		}
	}

	return fmt.Errorf(
		"%w: HiddenValueCallback not found",
		ErrAldiTalkCallbackPoWNoMatch,
	)
}
