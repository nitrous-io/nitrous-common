package sshkey

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
)

var ErrInvalidPublicKey = errors.New("sshkey: invalid public key")

var pubkeyRe = regexp.MustCompile(`[a-zA-Z0-9\-]+\s+([a-zA-Z0-9+\/]+=*)(\s|$)`)

func Md5Fingerprint(key string) (string, error) {
	m := pubkeyRe.FindStringSubmatch(key)
	if m == nil {
		return "", ErrInvalidPublicKey
	}

	data, err := base64.StdEncoding.DecodeString(m[1])
	if err != nil {
		return "", err
	}

	md5sum := md5.Sum(data)
	buf := &bytes.Buffer{}
	for _, v := range md5sum {
		buf.WriteString(fmt.Sprintf("%02x:", v))
	}

	return buf.String()[:47], nil
}
