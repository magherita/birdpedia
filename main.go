package main

import (
	"context"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/bird", getBirdHandler).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {
	connString := "mongodb://localhost:27017"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))

	if err != nil {
		panic(err)
	}

	defer client.Disconnect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	InitDatabase(&Connect{client: client})

	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// templates.ExecuteTemplate(w, "index.html", nil)
	// Use during development to avoid having to restart server
	// after every change in HTML
	t, _ := template.ParseFiles("./assets/index.html")
	t.Execute(w, nil)
}
