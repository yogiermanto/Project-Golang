package main

import (
	"fmt"
	"net/http"

	ctrl "project/controller"
)

type M map[string]interface{}

func main() {
	http.HandleFunc("/", ctrl.Index)
	http.HandleFunc("/index", ctrl.Index)
	http.HandleFunc("/cust-queue", ctrl.CustQueue)
	http.HandleFunc("/admin-queue", ctrl.AdminQueue)
	http.HandleFunc("/insert", ctrl.TambahAntrian)
	http.HandleFunc("/update", ctrl.UpdateMeja)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	fmt.Println("server started at localhost:4346")
	http.ListenAndServe(":4346", nil)
}
