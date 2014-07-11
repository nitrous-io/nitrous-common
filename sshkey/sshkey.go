package sshkey

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidPublicKey = errors.New("sshkey: invalid public key")

var pubkeyRe = regexp.MustCompile(
	`(?P<sigformat>[a-zA-Z0-9\-]+)` + // match 1: ssh-rsa, ssh-dss, etc
		`\s+` +
		`(?P<sigblob>[a-zA-Z0-9+\/]+=*)(\s|$)`) // match 2: actual key blob

func Md5Fingerprint(key string) (string, error) {
	m := pubkeyRe.FindStringSubmatch(key)
	if m == nil {
		return "", ErrInvalidPublicKey
	}

	validFormat := false
	for _, s := range []string{"ssh-rsa", "ssh-dss", "ecdsa"} {
		if strings.HasPrefix(m[1], s) {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return "", ErrInvalidPublicKey
	}

	data, err := base64.StdEncoding.DecodeString(m[2])
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
