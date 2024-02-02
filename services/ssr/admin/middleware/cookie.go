package middleware

import (
	"net/http"
)

type MiddlewareOpts func(*Middleware)

func NewMiddleware(next http.Handler, opts ...MiddlewareOpts) http.Handler {
	mw := Middleware{
		Next:     next,
		Secure:   true,
		HTTPOnly: true,
	}
	for _, opt := range opts {
		opt(&mw)
	}
	return mw
}

func WithSecure(secure bool) MiddlewareOpts {
	return func(m *Middleware) {
		m.Secure = secure
	}
}

func WithHTTPOnly(httpOnly bool) MiddlewareOpts {
	return func(m *Middleware) {
		m.HTTPOnly = httpOnly
	}
}

type Middleware struct {
	Next     http.Handler
	Secure   bool
	HTTPOnly bool
}

func ID(r *http.Request) (id string) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return
	}
	return cookie.Value
}

func (mw Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := ID(r)
	if id == "" {
		// ❗TODO
		// if redirecting you can set a cookie to pass a message to the /login page
		// login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	mw.Next.ServeHTTP(w, r)
}
