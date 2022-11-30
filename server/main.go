package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorm.io/gorm"
)

// File system folders that contain static resources/templates
var STATIC_DIR = "static"
var TEMPLATE_DIR = "templates"

func main() {
	db := ConfigureDb()

	// Serves our main template
	http.HandleFunc("/", handleRoot(db))
	// Serves our static css assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(STATIC_DIR))))
	// Handles POST requests to add new Movies to the database for display
	http.HandleFunc("/movie/", handleCreateMovie(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Starting server on port:", port)
	http.ListenAndServe(":"+port, nil)
}

func handleRoot(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			log.Print("Invalid path: ", r.URL.Path)
			handle404(w)
			return
		}

		// Query movies from the database
		var movies []Movie
		db.Find(&movies)

		data := struct {
			Movies []Movie
		}{
			Movies: movies,
		}

		renderTemplate(w, "movie_rating", data)
	}
}

func handleCreateMovie(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			handle404(w)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Print(err)
			redirect(w, r, "/")
			return
		}

		name := r.PostForm.Get("name")
		rating := r.PostForm.Get("rating")

		if name == "" {
			redirect(w, r, "/")
			return
		}

		rating_number, err := strconv.Atoi(rating)
		if err != nil || rating_number < 1 || rating_number > 5 {
			redirect(w, r, "/")
			return
		}

		db.Create(&Movie{Name: name, Rating: uint(rating_number)})

		// Redirect users back to the home page after creating their movie rating
		redirect(w, r, "/")
	}
}

func renderTemplate(w http.ResponseWriter, template_name string, data any) {
	// The base template is always included in whatever template we render.
	base_template := fmt.Sprintf("%v/base.tmpl", TEMPLATE_DIR)
	template_to_render := fmt.Sprintf("%v/%v.tmpl", TEMPLATE_DIR, template_name)
	t, err := template.ParseFiles(base_template, template_to_render)

	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found")
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}
