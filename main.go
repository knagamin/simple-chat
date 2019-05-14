package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/knagamin/simple-chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	if err := t.templ.Execute(w, data); err != nil {
		log.Fatal("problem occurred for executing template")
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "192.168.17.218")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// Setup Somniauth
	gomniauth.SetSecurityKey("ThisShou1dBeComplex!!")
	gomniauth.WithProviders(
		google.New(
			"865250558556-onf1eqrra6itt4evq72j47qj0aorlns6.apps.googleusercontent.com", // client ID
			"WxeMtPOVDf-qEYXpHkdWJKGX",                   // client seacret
			"http://localhost:8080/auth/callback/google", // callback URI
		),
	)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Web Server is starting... port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
