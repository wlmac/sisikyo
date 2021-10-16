package feed

import (
	"fmt"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestAnnToImageAndText(t *testing.T) {
	c := api.DefaultClient()
	resp := api.AnnResp{}
	err := c.Do(api.AnnReq{}, &resp)
	if err != nil {
		t.Fatalf("http: %s", err)
	}
	var ann api.Ann
	for _, ann2 := range resp {
		if ann2.Title == "National Day of Truth and Reconciliation" {
			ann = ann2
			break
		}
	}
	ann.XImageAlt, ann.XImageURL, _ = api.GetImageFromMd(ann.Body)
	height := 1000
	img, text, err := AnnToImageAndText(c, ann, image.Rect(0, 0, height, height))
	if err != nil {
		t.Fatalf("image: %s", err)
	}
	t.Log(text)
	file, err := os.Create("img.png")
	if err != nil {
		t.Fatalf("file: %s", err)
	}
	defer func(file *os.File) {
		err2 := file.Close()
		if err2 != nil {
			err = fmt.Errorf("file: %w", err2)
		}
	}(file)
	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
}
