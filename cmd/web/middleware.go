package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	// Adds cross-site scripting protection
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	// Logs each request. This is a method of the application struct,
	// but because it has the correct interface for a ServeHTTP it remains
	// valid. And now we also have access to other methods or data on the
	// application struct itself. Here we use the logger.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Using defer here just waits until all follow-on http.Handlers that
		// might run as middleware can do their thing. But if any of them panic
		// we can catch it just before this recoverPanic middleware runs, and
		// we gracefully handle it via the `recover` built-in function.
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection:close" header on the response
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
			// Note the '()' coming next - this is an anonymous func
			// that executes immediately.
		}()

		next.ServeHTTP(w, r)
	})
}
