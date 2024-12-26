package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	app := app.InitApp()
	app.SetHandler(mbr.Handler(web.Root))

	mbr.Dump()

	app.Run()
}

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TestMiddleware: %s\n", r.RequestURI)

		if ctx := mbr.Context(r); ctx != nil {
			fmt.Printf("mbr.Route %s => %s\n", ctx.Route().Name(), ctx.Route().FullPath())
		}

		next.ServeHTTP(w, r)
	})
}
