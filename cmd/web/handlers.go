package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, req *http.Request) {
	latestSnippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	data := app.newTemplateData(req)
	data.Snippets = latestSnippets

	app.render(writer, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(writer http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
			return
		}
		app.serverError(writer, err)
		return
	}

	data := app.newTemplateData(req)
	data.Snippet = snippet

	app.render(writer, http.StatusOK, "view.tmpl.html", data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreate(writer http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(writer, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, req *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot exceed 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal one of these three values: [1,7,365]")

	if !form.Valid() {
		data := app.newTemplateData(req)
		data.Form = form
		app.render(writer, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.sessionManager.Put(req.Context(), "flash", "Snippet successfully created!")

	http.Redirect(writer, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmailAddress(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddValidationError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmailAddress(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddGeneralError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	originURL := app.sessionManager.PopString(r.Context(), "originURL")
	if originURL != "" && originURL != "/user/logout" {
		http.Redirect(w, r, originURL, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) about(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	app.render(writer, http.StatusOK, "about.tmpl.html", data)
}

func (app *application) accountView(writer http.ResponseWriter, request *http.Request) {
	id := app.sessionManager.GetInt(request.Context(), "authenticatedUserID")

	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}
		app.serverError(writer, err)
		return
	}

	data := app.newTemplateData(request)
	data.User = user

	app.render(writer, http.StatusOK, "account.tmpl.html", data)
}

type accountPasswordUpdateForm struct {
	CurrentPassword     string `form:"currentPassword"`
	NewPassword         string `form:"newPassword"`
	ConfirmNewPassword  string `form:"confirmNewPassword"`
	validator.Validator `form:"-"`
}

func (app *application) accountPasswordUpdate(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = accountPasswordUpdateForm{}

	app.render(writer, http.StatusOK, "password.tmpl.html", data)
}

func (app *application) accountPasswordUpdatePost(writer http.ResponseWriter, request *http.Request) {
	var form accountPasswordUpdateForm

	err := app.decodePostForm(request, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.ConfirmNewPassword), "confirmNewPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(form.NewPassword == form.ConfirmNewPassword, "confirmNewPassword", "Password don't match")

	if !form.Valid() {
		data := app.newTemplateData(request)
		data.Form = form
		app.render(writer, http.StatusUnprocessableEntity, "password.tmpl.html", data)
		return
	}

	id := app.sessionManager.GetInt(request.Context(), "authenticatedUserID")

	err = app.users.PasswordUpdate(id, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddValidationError("currentPassword", "Incorrect password")
			data := app.newTemplateData(request)
			data.Form = form
			app.render(writer, http.StatusUnprocessableEntity, "password.tmpl.html", data)
			return
		}
		app.serverError(writer, err)
		return
	}

	app.sessionManager.Put(request.Context(), "flash", "Password changed successfully")
	http.Redirect(writer, request, "/account/view", http.StatusSeeOther)
}
