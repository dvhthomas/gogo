package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"dvhthomas/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	connStr       string
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
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

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	app.infoLog.Printf("Connected to the database as %s", *dbUser)
	defer db.Close()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
