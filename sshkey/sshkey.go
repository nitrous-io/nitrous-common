package sshkey

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

var ErrInvalidPublicKey = errors.New("sshkey: invalid public key")

func Md5Fingerprint(key string) (string, error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(key)))

	if err != nil {
		return "", ErrInvalidPublicKey
	}

	md5sum := md5.Sum(pk.Marshal())
	buf := &bytes.Buffer{}
	for _, v := range md5sum {
		buf.WriteString(fmt.Sprintf("%02x:", v))
	}

	return buf.String()[:47], nil
}
