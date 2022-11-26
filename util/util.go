package util

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"

	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/secure"
)

func SignCheck(secret, HEXnoncesign string) (e error) {
	noncesign, e := hex.DecodeString(HEXnoncesign)
	if e != nil {
		return ecode.ErrConfigDataBroken
	}
	if len(noncesign) < 8+64 {
		return ecode.ErrConfigDataBroken
	}
	oldsign := make([]byte, 64)
	copy(oldsign, noncesign[len(noncesign)-64:])
	newsign := sha512.Sum512(append(noncesign[:len(noncesign)-64], secret...))
	if !bytes.Equal(oldsign, newsign[:]) {
		return ecode.ErrWrongSecret
	}
	return nil
}
func SignMake(secret string, nonce []byte) (HEXnoncesign string) {
	tmp := make([]byte, 8+len(nonce))
	binary.BigEndian.PutUint64(tmp, uint64(len(nonce)))
	copy(tmp[8:], nonce)
	newsign := sha512.Sum512(append(tmp, secret...))
	return hex.EncodeToString(append(tmp, newsign[:]...))
}
func Encrypt(secret string, plaintext []byte) (HEXciphertext string, e error) {
	ciphertext, e := secure.AesEncrypt(secret, plaintext)
	if e != nil {
		return "", ecode.ErrWrongSecret
	}
	return hex.EncodeToString(ciphertext), nil
}
func Decrypt(secret, HEXciphertext string) (plaintext []byte, e error) {
	ciphertext, e := hex.DecodeString(HEXciphertext)
	if e != nil {
		return nil, ecode.ErrConfigDataBroken
	}
	plaintext, e = secure.AesDecrypt(secret, ciphertext)
	if e == secure.ErrAesSecretLength || e == secure.ErrAesSecretWrong {
		e = ecode.ErrWrongSecret
	} else if e == secure.ErrAesCipherTextBroken {
		e = ecode.ErrConfigDataBroken
	}
	return
}
