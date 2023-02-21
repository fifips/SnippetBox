package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"time"
)

type application struct {
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	debugMode      bool
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templates      TemplateCache
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return db, nil
	}

	return db, nil
}

func main() {
	serverAddress := flag.String("addr", "localhost", "HTTP network address")
	serverPort := flag.Int("port", 4000, "HTTP network port")
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := openDb(*dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	app := application{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		debugMode:      *debug,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templates:      templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: scs.New(),
	}

	app.sessionManager.Store = mysqlstore.New(db)
	app.sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager.Cookie.Secure = true

	srv := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", *serverAddress, *serverPort),
		Handler:  app.routes(),
		ErrorLog: app.errorLogger,
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLogger.Printf("Starting server on %s:%d", *serverAddress, *serverPort)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.errorLogger.Fatal(err)
}
