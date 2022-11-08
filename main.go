package main

import (
	"net/http"
	"time"

	v "paket/valechotter"

	"github.com/labstack/echo/v4"
)

type (
	Validator func(interface{}) error
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		var req struct {
			Items   []string `json:"items" query:"items"`
			Name    *bool    `json:"name"`
			Friends []string `json:"friends"`
		}
		c.Bind(&req)
		return c.JSON(http.StatusOK, req)
	}, v.EchoValidator().
		Query(v.M{
			"items": v.String,
		}).
		Body(v.Object(v.M{
			"name":    v.Nullable(v.Bool),
			"friends": v.Array(v.String),
			"attr": v.Object(v.M{
				"width":  v.UuidV4,
				"height": v.UuidV4,
			}),
			"birth": v.Date().After(time.Now().Add(-24 * time.Hour)).Build(),
		})).Middleware())

	e.Logger.Fatal(e.Start(":1881"))
}
