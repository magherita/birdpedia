package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Bird ...
type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	/*
		The list of birds is now taken from the store instead of the package level  `birds` variable we had earlier

		The `store` variable is the package level variable that we defined in
		`store.go`, and is initialized during the initialization phase of the
		application
	*/
	birds, err := db.GetBirds()

	// Everything else is the same as before
	birdListBytes, err := json.Marshal(birds)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	bird := Bird{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bird.Species = r.Form.Get("species")
	bird.Description = r.Form.Get("description")

	// The only change we made here is to use the `CreateBird` method instead of
	// appending to the `bird` variable like we did earlier
	err = db.CreateBird(&bird)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
