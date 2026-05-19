package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"snippetbox.abdulmoiz.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
)
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string] *template.Template
	formDecoder *form.Decoder
}


func main() {
	addr := flag.String("addr", ":4000", "HTTP Network address block")
	dsn := flag.String("dsn", "web:abc@/shirearchive?parseTime=true", "MySQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache, err := newTemplateCache()
	defer db.Close()
	formDecoder := form.NewDecoder()
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		
		
	}
	
	// A very important pattern when we have to pass in multiple dependencies to handler is the following
	// make a struct containg all the dependencies in main
	//call mux.Handle() with function that you defined which returns a handler itself i.e function returns another function
	//also called closure
	// The reason we dont diretly define handler itself is that handler only takes 2 arguements
	// by defining a new function and passing our config struct we can use the dependencies in the function we return
	// Important design pattern
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
func openDB(dsn string)(*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

