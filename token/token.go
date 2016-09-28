package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

type Token struct {
	Raw       string `json:"raw,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Used      bool   `json:"used,omitempty"`
}

func NewToken(raw, sub string, exp int64) *Token {
	return &Token{Raw: raw, Subject: sub, ExpiresAt: exp}
}

func (t *Token) SigningString() (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	buf := [][]byte{[]byte(t.Raw), b, []byte(strconv.FormatInt(t.ExpiresAt, 10))}
	encbuf := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		encbuf[i] = base64Encode(buf[i])
	}
	return string(bytes.Join(encbuf, []byte("."))), nil
}

func Verify(signature string) bool {
	t, err := Decode(signature)
	if err != nil {
		return false
	}
	err = t.Valid()
	if err != nil {
		return false
	}
	return true
}

func Decode(signature string) (*Token, error) {
	decSli := strings.Split(signature, ".")
	if len(decSli) != 3 {
		return nil, fmt.Errorf("it is not a valid token, required 3 fields")
	}
	slice := make([][]byte, 3)
	var err error
	for i := 0; i < 3; i++ {
		slice[i], err = base64Decode([]byte(decSli[i]))
		if err != nil {
			return nil, fmt.Errorf("error when decode string field %s: tracing(%s)", decSli[i], err.Error())
		}
	}
	var t Token
	err = json.Unmarshal(slice[1], &t)
	if err != nil {
		return nil, err
	}
	if string(slice[0]) != t.Raw {
		return nil, fmt.Errorf("field 0 must equal to t.Raw")
	}
	return &t, nil
}

func (t *Token) Valid() error {
	if t.Used {
		return fmt.Errorf("token has been used")
	}
	now := time.Now().Unix()
	if t.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(t.ExpiresAt, 0))
		return fmt.Errorf("token is expired by %v", delta)
	}
	return nil
}

func (t *Token) VerifyExpiresAt(cmp int64, req bool) bool {
	return varifyExp(t.ExpiresAt, cmp, req)
}

func varifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}
