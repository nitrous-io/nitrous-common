package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// use `base64.NewDecoder` to properly handle newline
func decodeBase64(s string) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(s))
	return ioutil.ReadAll(decoder)
}

func padSlice(mode cipher.BlockMode, src []byte) []byte {
	// src must be a multiple of block size
	bs := mode.BlockSize()
	mult := int((len(src) / bs) + 1)
	leng := bs * mult
	// make padded slice
	src_padded := make([]byte, leng)
	copy(src_padded, src)
	return src_padded
}

func Encrypt(key, data []byte) (iv []byte, ciphertext []byte, err error) {
	// create aes-256 cipher
	h := sha256.New()
	h.Write(key)
	block, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		return nil, nil, err
	}
	// initialize ciphertext and iv
	iv = make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}
	// encrypt ciphertext
	mode := cipher.NewCBCEncrypter(block, iv)
	if len(data)%mode.BlockSize() != 0 {
		data = padSlice(mode, data)
	}
	ciphertext = make([]byte, len(data))
	mode.CryptBlocks(ciphertext, data)
	return []byte(encodeBase64(iv)), []byte(encodeBase64(ciphertext)), nil
}

func Decrypt(key, iv, encrypted []byte) (data []byte, err error) {
	// create aes-256 cipher
	h := sha256.New()
	h.Write(key)
	block, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		return nil, err
	}
	if iv, err = decodeBase64(string(iv)); err != nil {
		return nil, err
	}
	if encrypted, err = decodeBase64(string(encrypted)); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	if len(encrypted)%mode.BlockSize() != 0 {
		encrypted = padSlice(mode, encrypted)
	}
	mode.CryptBlocks(encrypted, encrypted)
	return encrypted, nil
}
