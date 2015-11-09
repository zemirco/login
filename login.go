package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/juju/errgo"
)

// HTTPError occurs when handling http requests
type HTTPError struct {
	Err     error
	Message string
	Code    int
}

func (e *HTTPError) Error() string {
	return e.Message
}

// Handler wraps custom http handler funcs
type Handler func(http.ResponseWriter, *http.Request) *HTTPError

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("request error: %s", errgo.Details(e.Err))
		http.Error(w, e.Message, e.Code)
	}
}

func InternalServerError(err error) *HTTPError {
	return &HTTPError{
		err,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	}
}

var (
	loginTemplate *template.Template
)

func init() {
	loginTemplate = template.Must(template.ParseFiles("layout.html", "login.html"))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Path("/login").Handler(Handler(GetLogin)).Methods("GET")
	// r.Path("/login").Handler(Handler(PostLogin)).Methods("POST")
	// r.Path("/logout").Handler(middleware.Restrict(Handler(PostLogout))).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))
}

// GetLogin handles GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) *HTTPError {
	if err := loginTemplate.Execute(w, nil); err != nil {
		return InternalServerError(err)
	}
	return nil
}
