package picker

import (
	"encoding/base64"
	fernet "github.com/fernet/fernet-go"
	"time"
	"net/http"
	"github.com/millisecond/olb/model"
)

const HTTP_COOKIENAME = "OLBKEY"

// Cookie contains TargetID
func EncryptCookie(key string, t *model.Target) (string, error) {
	k := fernet.MustDecodeKeys(key)
	plaintext := []byte(t.ID)
	tok, err := fernet.EncryptAndSign(plaintext, k[0])
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tok), nil
}

// decrypt from base64 to decrypted string
func DecryptCookie(key string, encrypted string) string {
	k := fernet.MustDecodeKeys(key)
	cyphertext := []byte(encrypted)
	msg := fernet.VerifyAndDecrypt(cyphertext, time.Hour, k)
	return string(msg)
}

// rndPicker picks a random target from the list of targets.
func CookieHTTPPicker(targets []*model.Target, w http.ResponseWriter, req *http.Request) *model.Target {
	cookie, _ := req.Cookie(HTTP_COOKIENAME)
	if cookie == nil {
		t := RndHTTPPicker(targets, w, req)
		cookie, _ := EncryptCookie("", t)
		http.SetCookie(w, &http.Cookie{Name: HTTP_COOKIENAME, Value: cookie, Path: "/"})
		return t
	} else {
		id := DecryptCookie("", cookie.Value)
		for _, t := range targets {
			if t.ID == id {
				return t
			}
		}
	}
	return RndHTTPPicker(targets, w, req)
}
