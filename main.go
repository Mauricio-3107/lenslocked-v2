package main

import (
	"fmt"
	"net/http"

	"github.com/Mauricio-3107/lenslocked-v2/controllers"
	"github.com/Mauricio-3107/lenslocked-v2/templates"
	"github.com/Mauricio-3107/lenslocked-v2/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	var usersC controllers.Users
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))

	r.Get("/signup", usersC.New)

	r.Get("/galleries", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "galleries.gohtml", "tailwind.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
