package image2

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"net/url"

	"github.com/disintegration/imaging"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/feed/text"
)

type tmplNamespace struct {
	Org      api.Org
	URL      string
	Author   api.UserResp
	Ann      api.Ann
	IconURL  string
	ImageURL string
}

func resizeCenter(size int, orig image.Image, filter imaging.ResampleFilter) (result *image.NRGBA) {
	result = imaging.Resize(orig, size, 0, filter)
	xImg2 := image.NewNRGBA(result.Bounds())
	result = imaging.Paste(xImg2, result, image.Point{
		X: 0,
		Y: size/2 - result.Bounds().Dy()/2,
	})
	return
}

// AnnToImageAndText2 generates an SVG and accompanying alt text from an Ann.
func AnnToImageAndText2(c *api.Client, ann api.Ann) (svg, text string, err error) {
	if ann.XImageURL != "" {
		var xImg image.Image
		xImg, _, err = getImage(c, ann.XImageURL)
		if err != nil {
			err = fmt.Errorf("xImg: %w", err)
			return
		}
		xImg = resizeCenter(1000, xImg, imaging.Lanczos)
		ann.XImageURL, err = dataURI(xImg)
		if err != nil {
			err = fmt.Errorf("xImg uri: %w", err)
			return
		}
	}

	ns := tmplNamespace{
		Ann: ann,
	}

	err = c.Do(api.UserReq{Username: string(ann.Author)}, &ns.Author)
	if err != nil {
		err = fmt.Errorf("author: %w", err)
		return
	}
	ns.Org, err = ann.ReqOrg(c)
	if err != nil {
		err = fmt.Errorf("org: %w", err)
		return
	}
	ns.URL = ann.URL(c).String()

	icon, _, err := getImage(c, ns.Org.Icon)
	if err != nil {
		err = fmt.Errorf("icon: %w", err)
		return
	}
	ns.IconURL, err = dataURI(icon)
	if err != nil {
		err = fmt.Errorf("icon uri: %w", err)
		return
	}

	var text2 string
	text2, err = body(ann)
	if err != nil {
		err = fmt.Errorf("body: %w", err)
		return
	}
	text += text2
	svg, err = svgTmpl(ns)
	if err != nil {
		err = fmt.Errorf("svg: %w", err)
		return
	}
	return
}

func svgTmpl(ns tmplNamespace) (body string, err error) {
	buf := new(bytes.Buffer)
	err = Tmpl.ExecuteTemplate(buf, "ann", ns)
	if err != nil {
		return
	}
	body = buf.String()
	return
}

func body(ann api.Ann) (body string, err error) {
	buf := new(bytes.Buffer)
	err = text.Tmpl.ExecuteTemplate(buf, "feed-body", ann)
	if err != nil {
		return
	}
	body = buf.String()
	return
}

func getImage(c *api.Client, u string) (img draw.Image, name string, err error) {
	if u[0] == '/' {
		var u2 *url.URL
		u2, err = url.Parse(u)
		if err != nil {
			err = fmt.Errorf("url: %w", err)
			return
		}
		u = c.BaseURL().ResolveReference(u2).String()
	}

	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}
	resp, err := c.HTTPClient().Do(request)
	if err != nil {
		err = fmt.Errorf("http: %w", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			err = fmt.Errorf("close: %w", err2)
		}
	}(resp.Body)
	img2, name, err := image.Decode(resp.Body)
	if err != nil {
		err = fmt.Errorf("decode: %w", err)
		return
	}
	img = image.NewRGBA(img2.Bounds())
	draw.Draw(img, img2.Bounds(), img2, image.Point{}, draw.Over)
	return
}

func dataURI(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return "", err
	}
	return DataURI("image/png", buf.Bytes()), nil
}

// DataURI makes a data: URI.
func DataURI(mimeType string, src []byte) string {
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(src)
}
