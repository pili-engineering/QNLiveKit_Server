package seekable

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeekable_EOFIfReqAlreadyParsed(t *testing.T) {
	body := "a=1"
	req, err := http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "3")
	req.ParseForm()
	_, err = New(req)
	assert.Equal(t, err.Error(), "EOF")
}

func TestSeekable_WorkaroundForEOF(t *testing.T) {
	body := "a=1"
	req, err := http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "3")
	_, _ = New(req)
	req.ParseForm()
	assert.Equal(t, req.FormValue("a"), "1")
	_, err = New(req)
	assert.NoError(t, err)
}

func testSeekabler(t *testing.T, r io.ReadCloser, data string) {
	b, _ := ioutil.ReadAll(r)
	assert.Equal(t, data, string(b))
	sr, ok := r.(Seekabler)
	assert.True(t, ok)
	err := sr.SeekToBegin()
	assert.NoError(t, err)
	b, _ = ioutil.ReadAll(sr)
	assert.Equal(t, data, string(b))
}

func TestSeekable(t *testing.T) {
	body := "a=1"
	req, err := http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "3")
	_, err = New(req)
	assert.NoError(t, err)
	testSeekabler(t, req.Body, body)

	req, err = http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = New(req)
	assert.NoError(t, err)
	testSeekabler(t, req.Body, body)
}

func TestSeekableLength(t *testing.T) {
	old := MaxBodyLength
	defer func() {
		MaxBodyLength = old
	}()
	MaxBodyLength = 2
	body := "a=1"
	req, err := http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", "3")
	_, err = New(req)
	assert.Equal(t, ErrTooLargeBody, err)
	b, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, body, string(b))

	body = "a=11111111"
	req, err = http.NewRequest("POST", "/a", bytes.NewBufferString(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ContentLength = -1
	_, err = New(req)
	assert.Equal(t, ErrTooLargeBody, err)
	b, _ = ioutil.ReadAll(req.Body)
	assert.Equal(t, body, string(b))
}
