package rest

import (
	"context"
	"errors"
	"golang.org/x/net/http2"
	"html/template"
	"mini-network/internal/v1/models"
	"mini-network/internal/v1/transport"
	"mini-network/internal/v1/usecase"
	errors2 "mini-network/pkg/errors"
	"net/http"
	"time"
)

type HttpTransport struct {
	server  *http.Server
	useCase usecase.User
}

func NewHttpTransport(useCase usecase.User) *HttpTransport {
	return &HttpTransport{useCase: useCase}
}

func (s *HttpTransport) Run(port string) error {
	mux := http.NewServeMux()
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	http2.ConfigureServer(s.server, &http2.Server{})
	s.initHanlders(mux)
	return s.server.ListenAndServeTLS("./cert/server.crt", "./cert/server.key")
}

func (s *HttpTransport) initHanlders(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("GET /login/{$}", s.loginForm)
	mux.HandleFunc("GET /registration/{$}", s.registrationForm)
	mux.HandleFunc("POST /login/{$}", s.login)
	mux.HandleFunc("POST /registration/{$}", s.registration)
	mux.Handle("/", transport.Middleware(s.index, s.useCase.GetTokenService()))
}

func (s *HttpTransport) index(w http.ResponseWriter, r *http.Request, user models.User) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	if err := tmpl.Execute(w, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HttpTransport) loginForm(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Valid() == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/login.html"))
	data := struct {
		Greeting string
		Index    string
	}{
		Greeting: "Hello",
		Index:    "Joe",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HttpTransport) registrationForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/registration.html"))
	data := struct {
		Greeting string
		Index    string
	}{
		Greeting: "Hello",
		Index:    "Joe",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HttpTransport) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	formData := r.Form
	email := formData.Get("email")
	password := formData.Get("password")

	if email == "" || password == "" {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	token, err := s.useCase.LoginUser(r.Context(), email, password)
	if err != nil {
		var errCode errors2.CodedError
		if errors.As(err, &errCode) {
			http.Error(w, errCode.Error(), errCode.GetCode())
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(s.useCase.GetTokenService().GetExpired()),
	})
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (s *HttpTransport) registration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	formData := r.Form
	name := formData.Get("name")
	email := formData.Get("email")
	password := formData.Get("password")
	if name == "" || email == "" || password == "" {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	if _, err := s.useCase.CreateUser(r.Context(), name, email, password); err != nil {
		var errCode errors2.CodedError
		if errors.As(err, &errCode) {
			http.Error(w, errCode.Error(), errCode.GetCode())
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
}

func (s *HttpTransport) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
