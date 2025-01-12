package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aquilax/truncate"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (app *application) createDomainCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Allowed methods: POST")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request: %v", err)
		return
	}

	domain, err := app.domains.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error find domain: %v", err)
		return
	}

	resp, err := http.Get(domain.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error find domain: %v", err)
		return
	}

	defer resp.Body.Close()

	var (
		h1Text, keywordsText, descriptionText string
	)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error find domain: %v", err)
		return
	}

	if h1 := doc.Find("h1").First(); h1.Length() > 0 {
		h1Text = truncate.Truncate(h1.Text(), 30, "...", truncate.PositionEnd)
	}

	if keywords := doc.Find(`meta[name="keywords"]`).First(); keywords.Length() > 0 {
		keywordsText = truncate.Truncate(keywords.Text(), 30, "...", truncate.PositionEnd)
	}

	if description := doc.Find(`meta[name="description"]`).First(); description.Length() > 0 {
		descriptionText = truncate.Truncate(description.Text(), 30, "...", truncate.PositionEnd)
	}

	_, err = app.domainChecks.Insert(domain.Id, resp.StatusCode, h1Text, keywordsText, descriptionText)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error insert domain: %v", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d/", domain.Id), http.StatusSeeOther)
}

func (app *application) domainsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Allowed methods are: GET, POST")
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.showDomains(w, r)
	case http.MethodPost:
		app.createDomain(w, r)
	}
}

func (app *application) createDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Allowed methods: POST")
		return
	}

	name := r.PostFormValue("name")
	name = strings.TrimSpace(name)

	_, err := url.ParseRequestURI(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid URI")
		return
	}

	id, err := app.domains.Insert(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error inserting domain: %v", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d/", id), http.StatusSeeOther)
}

func (app *application) showDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Allowed methods: GET")
		return
	}

	files := []string{
		"./web/templates/domains/index.tmpl",
		"./web/templates/layouts/page.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't parse template")
		return
	}

	domains, err := app.domains.GetAll()
	checks, err := app.domainChecks.GetAll()

	err = ts.Execute(w, map[string]any{
		"domains": domains,
		"checks":  checks,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't parse template")
		return
	}
}

func (app *application) showDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Allowed methods are: GET")
		return
	}

	files := []string{
		"./web/templates/domains/show.tmpl",
		"./web/templates/layouts/page.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't parse template")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't parse id %v", id)
		return
	}

	domain, err := app.domains.Get(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't find domain: %v", err)
		return
	}

	checks, _ := app.domainChecks.GetByDomain(domain.Id)

	err = ts.Execute(w, map[string]any{
		"domain": domain,
		"checks": checks,
	})
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./web/templates/home.tmpl",
		"./web/templates/layouts/page.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
