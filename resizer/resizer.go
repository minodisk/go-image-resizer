package resizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/nfnt/resize"
	"github.com/stretchr/goweb"
	"github.com/stretchr/goweb/context"
)

type Req struct {
	Uri   string `json:"uri"`
	Width uint   `json:"width"`
}

type Res struct {
	Uri string `json:"uri"`
}

type Err struct {
	Message string `json:"message"`
}

type ResizerController struct {
	Storage string
}

// ReadMany reads many people.
func (c *ResizerController) ReadMany(ctx context.Context) error {
	return goweb.Respond.With(ctx, 200, []byte("Post JSON data"))
}

// // Read reads one person.
// func (c *ResizerController) Read(id string, ctx context.Context) error {
//
// }

// Create creates a new person.
func (c *ResizerController) Create(ctx context.Context) error {
	body, err := ctx.RequestBody()
	if err != nil {
		log.Printf("%+v", err)
		e := Err{}
		e.Message = fmt.Sprintf("%q", err)
		buf, err := json.Marshal(e)
		if err != nil {
		}
		return goweb.Respond.With(ctx, http.StatusBadRequest, buf)
	}
	req := Req{}
	json.Unmarshal(body, &req)
	log.Printf("%+v", req)

	imageRaw, err := fetchImage(req.Uri)
	if err != nil {
		log.Printf("%+v", err)
		e := Err{}
		e.Message = fmt.Sprintf("%q", err)
		buf, err := json.Marshal(e)
		if err != nil {
		}
		return goweb.Respond.With(ctx, http.StatusBadRequest, buf)
	}
	imageResized := resize.Resize(req.Width, 0, imageRaw, resize.Lanczos3)

	out, err := os.Create(path.Join(c.Storage, "test.png"))
	if err != nil {
		log.Printf("%+v", err)
		e := Err{}
		e.Message = fmt.Sprintf("%q", err)
		buf, err := json.Marshal(e)
		if err != nil {
		}
		return goweb.Respond.With(ctx, http.StatusBadRequest, buf)
	}
	defer out.Close()
	png.Encode(out, imageResized)

	log.Printf("%+v", ctx.Path())
	return goweb.Respond.With(ctx, 200, []byte("success"))
}

func fetchImage(uri string) (image.Image, error) {
	res, err := http.Get(uri)
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("The image doesn't exist at the URL")
	}

	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, errors.New("Can't decode")
	}

	return img, nil
}
