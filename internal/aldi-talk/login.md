```go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// ---- Config ----
// Passe diese Werte bei Bedarf an.
const (
	base = "https://login.alditalk-kundenbetreuung.de"
	realmPath = "/signin/json/realms/root/realms/alditalk"

	authenticateEP = base + realmPath + "/authenticate"
	sessionsEP     = base + "/signin/json/sessions"
	usersEP        = base + realmPath + "/users"
)

// ForgeRock AM typische Felder
// Siehe: NameCallback / PasswordCallback u.ä.
type frResponse struct {
	AuthID    string          `json:"authId"`
	Stage     string          `json:"stage"`
	Callbacks json.RawMessage `json:"callbacks"`
	// Bei Erfolg kommt ggf. ein Token oder Success-URL, wir brauchen nur nächsten Schritt.
}

type frCallback struct {
	Type   string           `json:"type"`
	Input  []frInput        `json:"input"`
	Output []frOutput       `json:"output"`
	Prompt string           `json:"prompt,omitempty"` // fallback
	Raw    *json.RawMessage `json:"-"`
}

type frInput struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type frOutput struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Benutzer-/Kundendaten (Teilmenge) – die Struktur der /users Antwort kann variieren.
// Wir extrahieren robust eine plausible Customer-ID aus mehreren möglichen Feldern.
type usersResponse struct {
	Result json.RawMessage `json:"result"`
	UserID string          `json:"userId"`
	UID    string          `json:"uid"`
	ID     string          `json:"id"`
	// Manche Installationen schicken ein Array von Treffern.
}

func main() {
	var username string
	var password string

	c := colly.NewCollector(
		colly.AllowedDomains("login.alditalk-kundenbetreuung.de", "www.alditalk-kundenbetreuung.de", "www.alditalk-kundenportal.de"),
		colly.MaxDepth(3),
	)
	c.WithTransport(&http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  false,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	})
	c.SetRequestTimeout(30 * time.Second)

	// Cookies automatisch verwalten
	c.AllowURLRevisit = true
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("HTTP %d on %s: %v\nBody: %s\n", r.StatusCode, r.Request.URL, err, truncate(string(r.Body), 400))
	})

	// Schritt 1: Challenge anfordern
	authResp1, err := postJSON(c, authenticateEP, map[string]any{})
	if err != nil {
		log.Fatalf("Authenticate (Step 1) fehlgeschlagen: %v", err)
	}

	var fr1 frResponse
	if err := json.Unmarshal(authResp1, &fr1); err != nil {
		log.Fatalf("Kann Authenticate-Response nicht parsen: %v", err)
	}
	if fr1.AuthID == "" {
		log.Fatalf("Keine authId in erster Authenticate-Antwort – Server-Flow unerwartet.")
	}

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

	// Hinweis: Manche Installationen liefern nach Step 2 bereits success.
	// Wir fahren mit Sessions + Users fort (aus deinem Log ersichtlich).
	_, err = postJSON(c, sessionsEP, map[string]any{"_action": "create"})
	if err != nil {
		log.Fatalf("Sessions anlegen fehlgeschlagen: %v", err)
	}

	// Schritt 3: Benutzer ermitteln (/users). Teilweise erwartet der Endpoint Filter/Query.
	// Wir versuchen mehrere Varianten: (a) leerer POST, (b) Query mit Username.
	cid, err := tryUsers(c, username)
	if err != nil {
		log.Fatalf("CID nicht ermittelbar: %v", err)
	}

	fmt.Println("CID:", cid)
}

func tryUsers(c *colly.Collector, username string) (string, error) {
	// Variante A: nackter POST
	if body, err := postJSON(c, usersEP, map[string]any{}); err == nil {
		if cid := extractCID(body); cid != "" {
			return cid, nil
		}
	}
	// Variante B: Query mit username (ForgeRock Directory-Style)
	u := usersEP + "?" + url.Values{"_queryFilter": {fmt.Sprintf("uid+eq+\"%s\"", escapeForQuery(username))}}.Encode()
	body, err := get(c, u)
	if err != nil {
		return "", err
	}
	if cid := extractCID(body); cid != "" {
		return cid, nil
	}
	return "", errors.New("users-Response enthielt keine erkennbare CID")
}

func fillCallbacks(fr frResponse, user, pass string) ([]map[string]any, error) {
	var cbArr []frCallback
	if err := json.Unmarshal(fr.Callbacks, &cbArr); err != nil {
		return nil, err
	}

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

func extractCID(body []byte) string {
	// Versuche mehrere Strukturen robust zu lesen.
	// 1) Einzelnes Objekt mit userId/uid/id
	var u usersResponse
	if err := json.Unmarshal(body, &u); err == nil {
		if u.UserID != "" {
			return u.UserID
		}
		if u.UID != "" {
			return u.UID
		}
		if u.ID != "" {
			return u.ID
		}
		// 2) Array-Hülle oder generische Map
	}
	// 2) Array von Objekten
	var arr []map[string]any
	if err := json.Unmarshal(body, &arr); err == nil {
		for _, it := range arr {
			if cid := firstString(it, "userId", "uid", "id", "customerId", "customerID"); cid != "" {
				return cid
			}
		}
	}
	// 3) Generische Map mit "result" Feld
	var m map[string]any
	if err := json.Unmarshal(body, &m); err == nil {
		if cid := firstString(m, "userId", "uid", "id", "customerId", "customerID"); cid != "" {
			return cid
		}
		if r, ok := m["result"].([]any); ok {
			for _, it := range r {
				if mm, ok := it.(map[string]any); ok {
					if cid := firstString(mm, "userId", "uid", "id", "customerId", "customerID"); cid != "" {
						return cid
					}
				}
			}
		}
	}
	return ""
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok2 := v.(string); ok2 && s != "" {
				return s
			}
		}
	}
	return ""
}

func postJSON(c *colly.Collector, endpoint string, payload any) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	reqHeaders.Set("Accept", "application/json, text/plain, */*")

	var respBody []byte
	var reqErr error

	err = c.Request("POST", endpoint, bytes.NewReader(b), nil, reqHeaders)
	if err != nil {
		return nil, err
	}

	// colly liefert Response im OnResponse-Hook; wir verwenden eine temporäre Collector-Kopie
	tmp := c.Clone()
	tmp.OnResponse(func(r *colly.Response) {
		if sameURL(r.Request.URL.String(), endpoint) && r.Request.Method == "POST" {
			respBody = r.Body
			if r.StatusCode >= 400 {
				reqErr = fmt.Errorf("HTTP %d: %s", r.StatusCode, truncate(string(r.Body), 200))
			}
		}
	})
	// Workaround: direkte GET um den Response-Hook zu triggern ist nicht nötig.
	// Stattdessen führen wir die eigentliche POST in tmp aus:
	err = tmp.Request("POST", endpoint, bytes.NewReader(b), nil, reqHeaders)
	if err != nil {
		return nil, err
	}
	if reqErr != nil {
		return nil, reqErr
	}
	if len(respBody) == 0 {
		return nil, errors.New("leere Antwort vom Server (POST)")
	}
	return respBody, nil
}

func get(c *colly.Collector, u string) ([]byte, error) {
	var body []byte
	var firstErr error
	cc := c.Clone()
	cc.OnResponse(func(r *colly.Response) {
		if sameURL(r.Request.URL.String(), u) {
			body = r.Body
			if r.StatusCode >= 400 {
				firstErr = fmt.Errorf("HTTP %d: %s", r.StatusCode, truncate(string(r.Body), 200))
			}
		}
	})
	if err := cc.Visit(u); err != nil {
		return nil, err
	}
	if firstErr != nil {
		return nil, firstErr
	}
	if len(body) == 0 {
		return nil, errors.New("leere Antwort vom Server (GET)")
	}
	return body, nil
}

func sameURL(a, b string) bool {
	// Vergleicht ohne übermäßig streng zu sein
	aa := strings.TrimSuffix(a, "/")
	bb := strings.TrimSuffix(b, "/")
	return aa == bb
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func escapeForQuery(s string) string {
	// sehr einfache Escapes für _queryFilter
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
```