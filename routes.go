package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/domains/{id}/checks/", app.createDomainCheck)
	mux.HandleFunc("/domains/{id}/", app.showDomain)
	mux.HandleFunc("/domains/", app.domainsHandler)
	mux.HandleFunc("/", app.home)

	return mux
}
