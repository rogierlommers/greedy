package articles

import (
	"encoding/binary"
	"net/url"

	"github.com/sirupsen/logrus"
)

func getHostnameFromUrl(addedUrl string) (hostname string) {
	u, err := url.Parse(addedUrl)
	if err != nil {
		logrus.Errorf("error looking up hostname [url: %s] [err: %s]", addedUrl, err)
	}
	return u.Host
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
