package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"server/static"

	"gorm.io/gorm"
)

func main() {
	db := configureDb()

	// Serves our main template
	http.HandleFunc("/", handleRoot(db))
	// Serves our static css assets
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.FS(static.StyleAssets))))
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

		renderFiles(w, "main", data)
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

func renderFiles(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFS(static.TemplateAssets, "base.tmpl", fmt.Sprintf("pages/%v.tmpl", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Not Found")
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}
