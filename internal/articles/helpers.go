package articles

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

func getMD5Hash(text string) string {
	toHash := text + time.Now().String()
	hash := md5.Sum([]byte(toHash))
	return hex.EncodeToString(hash[:])
}

func getHostnameFromUrl(addedUrl string) (hostname string) {
	u, err := url.Parse(addedUrl)
	if err != nil {
		log.Error("error looking up hostname from url", "url", addedUrl, "message", err)
	}
	return u.Host
}
