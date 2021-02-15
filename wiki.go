// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

type Page struct {
    Title string
    Body  []byte
}

var cachedTemplates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

func (p *Page) save() error {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0700)
	}
    filename := "data/" + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := "data/" + title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fileName := tmpl + ".html"
	
	err := cachedTemplates.ExecuteTemplate(w, fileName, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
	if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	
    log.Fatal(http.ListenAndServe(":8080", nil))
}