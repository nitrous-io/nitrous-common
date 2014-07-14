package cipher

import (
	"reflect"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("secret key")
	data := []byte("super secret data")
	iv, encrypted, err := Encrypt(key, data)
	if err != nil {
		t.Fatal("Failed to encrypt", err)
	}

	res, err := Decrypt(key, iv, encrypted)
	if reflect.DeepEqual(res, data) || err != nil {
		t.Fatal("Failed to decrypt", err)
	}
}
