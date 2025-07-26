package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

func MakeSECSecret(clientId, clientSecret string, t time.Time) string {
	utc := t.Format(http.TimeFormat)
	return fmt.Sprintf("SEC %x;%s",
		sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", clientId, clientSecret, utc))), utc)
}
