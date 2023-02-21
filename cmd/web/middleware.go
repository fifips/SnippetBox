package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		writer.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "deny")
		writer.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(writer, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		app.infoLogger.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(writer, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writer.Header().Set("Connection", "close")
				app.serverError(writer, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(writer, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.sessionManager.Put(request.Context(), "originURL", request.URL.Path)

		if !app.isAuthenticated(request) {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}
		writer.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(writer, request)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := app.sessionManager.GetInt(request.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(writer, request)
			return
		}

		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(writer, err)
			return
		}

		if exists {
			ctx := context.WithValue(request.Context(), isAuthenticatedContextKey, true)
			request = request.WithContext(ctx)
		}

		next.ServeHTTP(writer, request)
	})
}

func CSRFToken(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return csrfHandler
}
