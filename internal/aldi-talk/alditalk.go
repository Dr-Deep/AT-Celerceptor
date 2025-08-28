package alditalk

import (
	"regexp"

	_ "github.com/PuerkitoBio/goquery"
)

var (
	customerIDRegexp = regexp.MustCompile(`^C-\d{10}$`) // C-0123456789
)

/*

- ate
* config?
* login session '/login'
	* cookie banner
	* cookies speichern
	* einloggen
* datenvolumen? (int+einheit)
* gotify (success,fail)
* das alles im intervall

config: tel, passwd, intervall duration
useragent setzen
cookies speichern
*/
