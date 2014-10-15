package main

import (
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"github.com/minodisk/go-image-resizer/resizer"
	"github.com/stretchr/goweb"
)

func main() {
	port := os.Getenv("PORT")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	storage := path.Join(usr.HomeDir, ".go-image-resizer")
	os.Mkdir(storage, 0777)
	goweb.MapStatic("/storage", storage)

	controller := resizer.ResizerController{}
	controller.Storage = storage
	if err := goweb.MapController(&controller); err != nil {
		log.Fatal("Can't map ResizerController: ", err)
	}
	log.Printf("Start listening on port %s", port)
	if err := http.ListenAndServe(":"+port, goweb.DefaultHttpHandler()); err != nil {
		log.Fatal("Can't listen and serve: ", err)
	}
}
