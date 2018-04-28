package main

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	conn "system/daemon/controllers"
)

func init() {
	RouteMap = map[string]func(http.ResponseWriter, *http.Request){
		"/server/": ServerRoute,
	}
}

func ServerRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("r.URL.Path", r.URL.Path)
	if parts := strings.Split(r.URL.Path, "/"); len(parts) >= 3 {
		if method := reflect.ValueOf(&conn.Server{conn.NewController(w, r, cfg)}).MethodByName(strings.Title(parts[2])); method.IsValid() {
			method.Call(nil)
		}
		return
	}
	conn.NewController(w, r, cfg).Error()
}
