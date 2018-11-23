/*
Package examples shows how to implement a basic CRUD for two data structures with the api2go server functionality.
To play with this example server you can run some of the following curl requests

In order to demonstrate dynamic baseurl handling for requests, apply the --header="REQUEST_URI:https://www.your.domain.example.com" parameter to any of the commands.

Create a new user:
	curl -X POST http://localhost:31415/v0/users -d '{"data" : {"type" : "users" , "attributes": {"user-name" : "marvin"}}}'


List users:
	curl -X GET http://localhost:31415/v0/users

List paginated users:
	curl -X GET 'http://localhost:31415/v0/users?page\[offset\]=0&page\[limit\]=2'
OR
	curl -X GET 'http://localhost:31415/v0/users?page\[number\]=1&page\[size\]=2'

Update:
	curl -vX PATCH http://localhost:31415/v0/users/1 -d '{ "data" : {"type" : "users", "id": "1", "attributes": {"user-name" : "better marvin"}}}'

Delete:
	curl -vX DELETE http://localhost:31415/v0/users/2

Create a chocolate with the name sweet
	curl -X POST http://localhost:31415/v0/chocolates -d '{"data" : {"type" : "chocolates" , "attributes": {"name" : "Ritter Sport", "taste": "Very Good"}}}'

Create a user with a sweet
	curl -X POST http://localhost:31415/v0/users -d '{"data" : {"type" : "users" , "attributes": {"user-name" : "marvin"}, "relationships": {"sweets": {"data": [{"type": "chocolates", "id": "1"}]}}}}'

List a users sweets
	curl -X GET http://localhost:31415/v0/users/1/sweets

Replace a users sweets
	curl -X PATCH http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'

Add a sweet
	curl -X POST http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'

Remove a sweet
	curl -X DELETE http://localhost:31415/v0/users/1/relationships/sweets -d '{"data" : [{"type": "chocolates", "id": "2"}]}'
*/
package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hi-lap/api2go-gorm-gin-crud-example/model"
	"github.com/hi-lap/api2go-gorm-gin-crud-example/resource"
	"github.com/hi-lap/api2go-gorm-gin-crud-example/storage"
	_ "github.com/lib/pq"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/routing"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

type ginRouter struct {
	router *gin.Engine
}

func (g ginRouter) Handler() http.Handler {
	return g.router
}

func (g ginRouter) Handle(protocol, route string, handler routing.HandlerFunc) {
	wrappedCallback := func(c *gin.Context) {
		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		handler(c.Writer, c.Request, params, c.Keys)
	}

	g.router.Handle(protocol, route, wrappedCallback)
}

//Gin creates a new api2go router to use with the gin framework
func Gin(g *gin.Engine) routing.Routeable {
	return &ginRouter{router: g}
}

func main() {
	r := gin.Default()
	api := api2go.NewAPIWithRouting(
		"api",
		api2go.NewStaticResolver("/"),
		Gin(r),
	)

	db, err := storage.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userStorage := storage.NewUserStorage(db)
	chocStorage := storage.NewChocolateStorage(db)
	api.AddResource(model.User{}, resource.UserResource{ChocStorage: chocStorage, UserStorage: userStorage})
	api.AddResource(model.Chocolate{}, resource.ChocolateResource{ChocStorage: chocStorage, UserStorage: userStorage})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Run(":3000") // listen and serve on 0.0.0.0:31415
}
