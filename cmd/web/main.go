package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"dvhthomas/snippetbox/pkg/models"
	"dvhthomas/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// Rather than using a brittle string all over, we define a type
// and create context keys on a per-key basis. Strongly typed.
// NOTE: This should also keep conflicts to other middle keys to
// a minimum.
type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	connStr  string
	// We define the interfaces inline so that both the real implementation
	// of mysql/snippets.go *and* the mock snippets implementations can
	// match the interface instead of putting a concrete implementation like
	// mysql.UserModel in here instead.
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
	// Same for users -- describe the interface and not the concrete implementation.
	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	dbUser := flag.String("dbuser", "", "Database user that application runs under")
	dbPass := flag.String("dbpass", "", "Database password for the application user")
	dbHost := flag.String("dbhost", "0.0.0.0", "Database host")
	addr := flag.String("addr", ":4000", "HTTP network address")
	// The secret is a random 32 character value used to encrypt and auth cookies
	secret := flag.String("secret", "", "Secret key for session encryption.\nTry 'openssl rand -base64 32' to generate one")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// DSN is a Data Source Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/snippetbox?parseTime=true", *dbUser, *dbPass, *dbHost)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Sessions will always expire after 12 hours
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app.infoLog.Printf("Connected to the database as %s", *dbUser)
	defer db.Close()

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
