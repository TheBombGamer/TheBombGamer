package main

import (
    "fmt"
    "html/template"
    "net/http"
)

// Page data structure
type Page struct {
    Title  string
    Header string
}

// Function to render templates
func renderTemplate(w http.ResponseWriter, tmpl string, page Page) {
    t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.ExecuteTemplate(w, "layout", page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Home handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
    page := Page{Title: "Home", Header: "Welcome to My Go Web Server"}
    renderTemplate(w, "index.html", page)
}

// About handler
func aboutHandler(w http.ResponseWriter, r *http.Request) {
    page := Page{Title: "About", Header: "About Us"}
    renderTemplate(w, "about.html", page)
}

// Contact handler
func contactHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        name := r.FormValue("name")
        message := r.FormValue("message")
        fmt.Fprintf(w, "Thank you, %s! Your message has been received: %s", name, message)
        return
    }
    page := Page{Title: "Contact", Header: "Contact Us"}
    renderTemplate(w, "contact.html", page)
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func main() {
    // Serve static files from the "public" directory
    fs := http.FileServer(http.Dir("public"))
    http.Handle("/public/", http.StripPrefix("/public/", fs))

    // Register routes
    http.Handle("/", loggingMiddleware(http.HandlerFunc(homeHandler)))
    http.Handle("/about", loggingMiddleware(http.HandlerFunc(aboutHandler)))
    http.Handle("/contact", loggingMiddleware(http.HandlerFunc(contactHandler)))

    // Start the server on port 8080
    fmt.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
