package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/to-do/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createTodoForm))
	mux.Post("/to-do/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createTodo))
	mux.Get("/to-do/:id", dynamicMiddleware.ThenFunc(app.showTodo))

	mux.Post("/to-do/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateTodo))
	mux.Get("/to-do/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateTodoForm))

	mux.Post("/to-do/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteTodo))
	mux.Get("/to-do/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteTodoForm))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
