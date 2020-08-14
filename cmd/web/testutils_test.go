package main

import (
	"dvhthomas/snippetbox/pkg/models/mock"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
)

func newTestApplication(t *testing.T) *application {
	// This is a relative path to the current working directory
	// that `go test` is being run in. So it will fail if you're not
	// in the project root dir.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("a"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
		templateCache: templateCache,
	}
}

// Define a custom testServer type which anonymously embeds an http.Server instance
type testServer struct {
	*httptest.Server
}

// Create newTestServer which returns our new testServer type configured
// for TLS
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// We want to test routes without having them automatically redirect to
	// another route. So we want to intercept any redirects.
	// But we don't want to lose any cookies that get set along the way. So
	// first we set up a cookie jar as part of our new server and it's Client....
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// ...and now that the cookies are handled correctly and we'll be able to inspect
	// them, we can ask our server to stop redirects and instead just send
	// us the response that was trying to ask for the redirect. Meaning: send us
	// the response we want to look at, not the one it was trying to redirect us to.
	// Clear as mud?! :-)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test serer, and returns the response
// status code, headers, and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
