package alditalk

import "github.com/gocolly/colly/v2"

/*
haben wir nach login, cookies?



COOKIE SETZEN: rememberMe: true
session id behalten
CustomerID: tef_customer_id
UserID: user_id
is_logged: muss true sein sonst kein login
*/

// + session
func Login(c *colly.Collector, username string, password string) (customerID string, err error) {

	/*
		neue verbindung:
			login seite
			einloggen
			cookies akzeptieren
			aus cookies CID lesen : tef_customer_id=C-0005150729; oder regexp für cookies

		alte verbindung:
			noch eingeloggt? neu oder alte CID
	*/

}
