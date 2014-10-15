package main

import (
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"github.com/minodisk/go-image-resizer/resizer"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
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

	goweb.Map(func(ctx context.Context) error {
		goweb.Respond.With(ctx, 200, []byte("POST /resizer {url:\"[リサイズする画像のURL]\", width:[リサイズするサイズ]}"))
		return nil
	})

	log.Printf("Start listening on port %s", port)
	if err := http.ListenAndServe(":"+port, goweb.DefaultHttpHandler()); err != nil {
		log.Fatal("Can't listen and serve: ", err)
	}
}
