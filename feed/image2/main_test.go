package image2

import (
	"fmt"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"os"
	"testing"
)

func TestAnnToImageAndText2(t *testing.T) {
	c := api.DefaultClient()
	resp := api.AnnResp{}
	err := c.Do(api.AnnReq{}, &resp)
	if err != nil {
		t.Fatalf("http: %s", err)
	}
	var ann api.Ann
	for _, ann2 := range resp {
		if ann2.Title == "Gardening Club Raffle" {
			ann = ann2
			break
		}
	}
	ann.XImageAlt, ann.XImageURL, _ = api.GetImageFromMd(ann.Body)
	svg, text, err := AnnToImageAndText2(c, ann)
	if err != nil {
		t.Fatalf("image: %s", err)
	}
	t.Log(text)
	file, err := os.Create("img.svg")
	if err != nil {
		t.Fatalf("file: %s", err)
	}
	defer func(file *os.File) {
		err2 := file.Close()
		if err2 != nil {
			err = fmt.Errorf("file: %w", err2)
		}
	}(file)
	_, err = file.WriteString(svg)
	if err != nil {
		err = fmt.Errorf("write: %w", err)
		return
	}
}
