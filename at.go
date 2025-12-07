package celerceptor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

func (c *Celerceptor) initAldiTalk() {
	c.Lock()

	/*

		c.tryLogin()
		c.tryUpdateDataVol()
	*/

	//! login
	tokenID := c.login()
	_ = tokenID

	uid, err := c.atframClient.GetFramUIDFromSession()
	if err != nil {
		panic(err)
	}

	c.logger.Info("UID:", fmt.Sprintf("%#v", uid))

	c.Unlock()
}

/*
 * AT FUNCS
 */

/*
* ticker um datenvolumen einzusehen (smart ticker) => (nachbuchen oder nix?)
* volumen merken
 */

// returns tokenID
func (c *Celerceptor) login() string {
	c.logger.Info("trying to login")

	successResp, err := c.atframClient.Login()
	if err != nil {
		c.logger.Fatal("Login", err.Error())
	}

	//! follow success url

	resp, err := c.atframClient.DoHttpRequest(
		http.MethodGet,
		successResp.SuccessURL,
		[]byte{},
	)
	if err != nil {
		panic(err)
	}
	_ = resp

	//!

	sessionInfo, err := c.atframClient.GetFramSessionInfo()
	if err != nil {
		c.logger.Fatal("sessioninfo", err.Error())
	}

	userInfo, err := c.atframClient.GetFramUserInfo(sessionInfo.UserID)
	if err != nil {
		c.logger.Fatal("userinfo", err.Error())
	}

	c.logger.Info(
		"logged in as",
		color.BgRGB(27, 50, 130).Sprintf("+")+color.BgRGB(239, 124, 28).Sprintf("%s", userInfo.TelephoneNumber[0]),
	)

	return ""
}

func (c *Celerceptor) tryLogout() {
	c.logger.Info("trying to logout")
}

func (c *Celerceptor) tryUpdateDataVol() {
	c.logger.Info("trying to update datavol")

	//
}
func (c *Celerceptor) tryGetDataVol() {
	c.logger.Info("trying to get datavol")
}

// eig nur für debug logs
func _readToBuf(r io.ReadCloser) *bytes.Buffer {

	out, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	//? io.LimitReader()

	return bytes.NewBuffer(out)
}
