package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

type Req struct {
	Uri   string `json:"uri"`
	Width uint   `json:"width"`
}

func main() {
	log.Println("main")
	port := os.Getenv("PORT")
	http.HandleFunc("/", IndexHandler)
	log.Printf("Start listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	dec := json.NewDecoder(buf)
	var req Req
	dec.Decode(&req)
	fmt.Printf("%+v\n", req)

	img := fetchImage(req.Uri)
	m := resize.Resize(req.Width, 0, img, resize.Lanczos3)

	out, err := os.Create("storage/test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	png.Encode(out, m)

	fmt.Fprintln(w, "running")
}

func fetchImage(uri string) image.Image {
	res, err := http.Get(uri)
	if err != nil || res.StatusCode != 200 {
		log.Fatal("Doesn't exist")
	}

	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatal("Can't decode")
	}

	return img
}
