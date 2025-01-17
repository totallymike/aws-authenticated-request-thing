package request

import (
	"testing"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func newRequest() (signedRequest *SignedRequest) {
	signedRequest, _ = NewSignedRequest(
		"GET",
		"http://www.example.com/v1/network-hosts?foo=bar&baz=foo",
	)
	return
}

func TestCanonicalURI(t *testing.T) {
	t.Parallel()
	req := newRequest()

	if req.CanonicalURI != "/v1/network-hosts" {
		t.Errorf("%s != %s\n", req.CanonicalURI, "/v1/network-hosts")
	}
}

func TestCanonicalQueryString(t *testing.T) {
	t.Parallel()
	req := newRequest()

	expected := "baz=foo&\nfoo=bar"
	if req.CanonicalQueryString() != expected {
		t.Errorf("%s != %s\n", req.CanonicalQueryString(), expected)
	}
}

func TestCanonicalHeaders(t *testing.T) {
	t.Parallel()
	req := newRequest()

	expected := "content-type:application/vnd.api+json\n" +
		"host:www.example.com\n" +
		"x-amz-date:" + time.Now().UTC().Format("20060102T150405Z") + "\n"

	actual := req.CanonicalHeaders()
	if actual != expected {
		t.Errorf("%s != %s\n", actual, expected)
	}
}

func TestCanonicalHeadersWithSpaces(t *testing.T) {
	t.Parallel()
	req := newRequest()
	req.AddHeader("foo", `"   oh  yeah"`)
	req.AddHeader("bar", "   oh     yeah")

	expected := "bar:oh yeah\ncontent-type:application/vnd.api+json\n" +
		"foo:\"   oh  yeah\"\nhost:www.example.com\n" +
		"x-amz-date:" + time.Now().UTC().Format("20060102T150405Z") + "\n"


	if req.CanonicalHeaders() != expected {
		t.Errorf("\n%s != %s\n", req.CanonicalHeaders(), expected)
	}
}

func TestSignedHeaders(t *testing.T) {
	t.Parallel()

	req := newRequest()

	expected := "content-type;host;x-amz-date"
	actual := req.SignedHeaders()

	if expected != actual {
		t.Errorf("%s != %s\n", actual, expected)
	}
}

func TestSignedPayload(t *testing.T) {
	t.Parallel()

	req := newRequest()

	hash := sha256.New()
	io.WriteString(hash, "")
	expected := hex.EncodeToString(hash.Sum(nil))
	actual := req.SignedPayload("")

	if expected != actual {
		t.Errorf("%s != %s\n", actual, expected)
	}
}

func TestCanonicalRequest(t *testing.T) {
	t.Parallel()

	req := newRequest()

	signedPayload := signPayload("")
	expected := "GET\n" +
		"/v1/network-hosts\n" +
		"baz=foo&\nfoo=bar\n" +
		"content-type:application/vnd.api+json\n" +
		"host:www.example.com\n" +
		"x-amz-date:" + time.Now().UTC().Format("20060102T150405Z") + "\n" + "\n" +
		"content-type;host;x-amz-date\n" +
		signedPayload

	actual := req.CanonicalRequest("")
	if expected != actual {
		t.Errorf("%s != %s\n", actual, expected)
	}
}

func TestHashedCanonicalRequest(t *testing.T) {
	// This is a bit of a dummy, as both the
	// test and the actual code use the same
	// signing method

	t.Parallel()

	req := newRequest()

	canonicalRequest := req.CanonicalRequest("")
	expected := signPayload(canonicalRequest)

	actual := req.HashedCanonicalRequest("")

	if expected != actual {
		t.Errorf("%s != %s\n", actual, expected)
	}
}
