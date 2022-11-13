package main

import (
	"database/sql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
)

func routerWrapper(w http.ResponseWriter, r *http.Request,
	router http.HandlerFunc) {

	log.Println(r.RequestURI)
	router(w, r)
}

type user struct {
	id       int
	username string
}

func QryUsersById(id int) (*user, error) {
	db, err := sql.Open("sqlite3", "web-app-example.db")
	var u *user

	if err != nil {
		log.Println(err)
		return u, err
	}

	defer db.Close()
	row := db.QueryRow("SELECT id,username FROM users WHERE id=?", id)
	if row.Err() != nil {
		log.Println(row.Err().Error())
		return u, err
	}
	u = new(user)
	row.Scan(&u.id, &u.username)
	return u, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := iris.New()
	app.WrapRouter(routerWrapper)

	app.HandleDir("/static/", "./static/")
	app.RegisterView(iris.HTML("./views", ".html").Layout("layout.html"))

	app.Get("/", func(ctx iris.Context) {
		ctx.View("home.html")
	})

	app.Get("/template", func(ctx iris.Context) {
		ctx.ViewData("title", "Template page")
		id, err := ctx.URLParamInt("id")
		if err != nil {
			id = 1
		}
		u, err := QryUsersById(id)
		if err != nil {
			log.Println(err)
			return
		}
		ctx.ViewData("user", u.username)
		ctx.View("template.html")
	})

	app.Get("/about", func(ctx iris.Context) {
		ctx.View("about.html", nil)

	})

	app.Post("/contact-us", func(ctx iris.Context) {
		body, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(body))
		ctx.JSON(context.Map{"message": "稍后会有工作人员联系您！"})
	})

	// http://localhost:8080
	app.Listen(":8080")
}
