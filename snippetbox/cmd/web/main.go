package main

import (
	"database/sql" // New import
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // New import
	"snippetbox.alexedwards.net/internal/models"
)
type application struct {
    errorLog 		*log.Logger
    infoLog  		*log.Logger
    snippets 		*models.SnippetModel
    users    		*models.UserModel
    templateCache	map[string]*template.Template
    formDecoder		*form.Decoder
    sessionManager	*scs.sessionManager
    // db *sql.DB  // if you want to store your database connection
}


func main() {
addr := flag.String("addr", ":4000", "HTTP network address")
// Define a new command-line flag for the MySQL DSN string.
dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQLdata source name")

flag.Parse()
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t",
log.Ldate|log.Ltime|log.Lshortfile)
dsn := os.Getenv("DB_DSN")
if dsn == "" {
	errorLog.Fatal("DB_DSN environment variable not set")
}
db, err := openDB(dsn)
if err != nil {
errorLog.Fatal(err)
}
// We also defer a call to db.Close(), so that the connection pool is
//closed
// before the main() function exits.
defer db.Close()
app := &application{
errorLog: 		errorLog,
infoLog:		infoLog,
snippets:		&models.SnippetModel{DB:db},
users: 			&models.UsersModel{DB: db},
templateCache: 	templateCache,
formDecoder: 	formDecoder,
sessionManager: sessionManager,
}

tlsConfig := tls.Config{
	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
}

srv := &http.Server {
Addr: *addr,
ErrorLog: errorLog,
Handler: app.routes(),
TLSConfig: tlsConfig,
IdleTimeout:  time.Minute,
ReadTimeout:  5 * time.Second,
WriteTimeout: 10 * time.Second,
}

infoLog.Printf("Starting server on %s", *addr)
err = srv.ListenAndServe()
errorLog.Fatal(err)
}
// The openDB() function wraps sql.Open() and returns a sql.DB connection
//pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
