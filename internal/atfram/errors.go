package atfram

import "errors"

var (
	// AldiTalkConfig
	ErrAldiTalkConfigMissingBaseURI  = errors.New("we need a baseURI")
	ErrAldiTalkConfigMissingUsername = errors.New("we need a username")
	ErrAldiTalkConfigMissingPassword = errors.New("we need a password")

	// AldiTalkClient
	ErrAldiTalkClientInvalidStatusCode = errors.New("unexpected response status-code")
)
