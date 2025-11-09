package atfram

import "errors"

var (
	// AldiTalkConfig
	ErrAldiTalkConfigMissingBaseURI  = errors.New("we need a baseURI")
	ErrAldiTalkConfigMissingUsername = errors.New("we need a username")
	ErrAldiTalkConfigMissingPassword = errors.New("we need a password")

	// AldiTalkClient
	ErrAldiTalkClientInvalidStatusCode = errors.New("unexpected response status-code")
	ErrAldiTalkClientUnknownCallback   = errors.New("please use github issues, I do not know that callback yet or havent implemented it")

	// Callback Errors
	ErrAldiTalkCallbackEmptyPoWScript = errors.New("received empty proof of work script message?")
	ErrAldiTalkCallbackPoWNoMatch     = errors.New("cannot find proof of work vars")
	ErrAldiTalkCallbackPoW            = errors.New("hash does not satisfy difficutly")
)
