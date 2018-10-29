// This code from official go tutorial
// and I am using it only to understand basic concepts
package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templatePath = "templates/"
var filePath = "files/"

// Load and cache all our templates
var templates = template.Must(template.ParseFiles(templatePath+"edit.html", templatePath+"view.html"))

// Validation rules for url
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// This code describes our data structure
type Page struct {
	Title string
	Body  []byte
}

// This part of code add save method to our data structure
// this method return error if we have an issues with file creating
func (p *Page) save() error {
	filename := p.Title + ".txt"
	// 0600 - it's write read permissions
	return ioutil.WriteFile(filePath+filename, p.Body, 0600)
}

// Function reads our file data and returns Data object or error
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filePath + filename)
	if err != nil {
		return nil, err
	}
	// we are returning data object strucutre link
	return &Page{Title: title, Body: body}, nil
}

// Function validate our url path
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

// Function rendering a templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// getting templates from a cache list
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handler to view our docs
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// handler to edit our file
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// handler to save our data
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Clusure function that validate our url and getting filename from url
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
