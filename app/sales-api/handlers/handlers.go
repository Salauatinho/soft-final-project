package handlers

import (
	"github.com/Salauatinho/soft-final-project/business/auth"
	"github.com/Salauatinho/soft-final-project/business/data/cafe"
	"github.com/Salauatinho/soft-final-project/business/data/staff"
	"github.com/Salauatinho/soft-final-project/business/data/user"
	"github.com/Salauatinho/soft-final-project/business/mid"
	"github.com/Salauatinho/soft-final-project/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	cg := CheckGroup{
		build: build,
		db:    db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveness", cg.liveness)

	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/users/:page/:rows", ug.query, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodGet, "/users/:id", ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/users/token/:kid", ug.token)
	app.Handle(http.MethodPost, "/users", ug.create, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id", ug.update, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id", ug.delete, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))

	cafeg := cafeGroup{
		cafe: cafe.New(log, db),
		auth: a,
	}
	app.Handle(http.MethodGet, "/cafe/:page/:rows", cafeg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/cafe/:id", cafeg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/cafe", cafeg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/cafe/:id", cafeg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/cafe/:id", cafeg.delete, mid.Authenticate(a))

	sg := staffGroup{
		staff: staff.New(log, db),
		auth:  a,
	}
	app.Handle(http.MethodGet, "/users/:page/:rows", sg.query, mid.Authenticate(a))
	app.Handle(http.MethodGet, "/users/:id", sg.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/users", sg.create, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/users/:id", sg.update, mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/users/:id", ug.delete, mid.Authenticate(a))

	return app
}
