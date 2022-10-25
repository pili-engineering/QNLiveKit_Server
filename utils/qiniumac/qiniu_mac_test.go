package qiniumac

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sk = []byte("secret_key")
	su = "su_info"
)

func Test_Sign(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path/to/api?param=value", nil)
	req.Header.Set("Content-Type", "application/json")

	act, err := SignRequest(sk, req)
	assert.NoError(t, err)

	h := hmac.New(sha1.New, sk)
	h.Write([]byte("GET /path/to/api?param=value\nHost: example.com\nContent-Type: application/json\n\n"))
	exp := h.Sum(nil)

	assert.Equal(t, exp, act)
}

func Test_SignWithXQiniu(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path/to/api?param=value", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Qiniu-Meta-App", "value")

	act, err := SignRequest(sk, req)
	assert.NoError(t, err)

	h := hmac.New(sha1.New, sk)
	h.Write([]byte("GET /path/to/api?param=value\nHost: example.com\nContent-Type: application/json" +
		"\nX-Qiniu-Meta-App: value" +
		"\n\n"))
	exp := h.Sum(nil)

	assert.Equal(t, exp, act)
}

func Test_SignAdmin(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path/to/api?param=value", nil)
	req.Header.Set("Content-Type", "application/json")

	act, err := SignAdminRequest(sk, req, su)
	assert.NoError(t, err)

	h := hmac.New(sha1.New, sk)
	h.Write([]byte("GET /path/to/api?param=value\nHost: example.com\nContent-Type: application/json" +
		"\nAuthorization: QiniuAdmin " + su +
		"\n\n"))
	exp := h.Sum(nil)

	assert.Equal(t, exp, act)
}

func Test_SignAdminWithXQiniu(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path/to/api?param=value", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Qiniu-Meta-App", "value")

	act, err := SignAdminRequest(sk, req, su)
	assert.NoError(t, err)

	h := hmac.New(sha1.New, sk)
	h.Write([]byte("GET /path/to/api?param=value\nHost: example.com\nContent-Type: application/json" +
		"\nAuthorization: QiniuAdmin " + su +
		"\nX-Qiniu-Meta-App: value" +
		"\n\n"))
	exp := h.Sum(nil)

	assert.Equal(t, exp, act)
}

func Test_signQiniuHeaderValues(t *testing.T) {

	w := bytes.NewBuffer(nil)

	header := make(http.Header)
	header.Set("X-Qbox-Meta", "value")

	signQiniuHeaderValues(header, w)
	assert.Empty(t, w.String())

	header.Set("X-Qiniu-Cxxxx", "valuec")
	header.Set("X-Qiniu-Bxxxx", "valueb")
	header.Set("X-Qiniu-axxxx", "valuea")
	header.Set("X-Qiniu-e", "value")
	header.Set("X-Qiniu-", "value")
	header.Set("X-Qiniu", "value")
	header.Set("", "value")
	signQiniuHeaderValues(header, w)

	assert.Equal(t, `
X-Qiniu-Axxxx: valuea
X-Qiniu-Bxxxx: valueb
X-Qiniu-Cxxxx: valuec
X-Qiniu-E: value`, w.String())
}
