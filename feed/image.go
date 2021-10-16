package feed

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/feed/font"
	"gitlab.com/mirukakoro/sisikyo/feed/text"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
)

const ratioOfIconToImage = 0.2

var resizeFilter = imaging.NearestNeighbor

func AnnToImageAndText(c *api.Client, ann api.Ann, rect image.Rectangle) (img *image.NRGBA, text string, err error) {
	img = image.NewNRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: image.Black}, image.Point{}, draw.Src)
	if ann.XImageURL != "" {
		err = drawImage(img, c.HTTPClient(), ann.XImageURL, func(dst, src image.Image) {
			srcDx := src.Bounds().Dx()
			srcDy := src.Bounds().Dy()
			srcRatio := float64(srcDx) / float64(srcDy)
			dstDx := dst.Bounds().Dx()
			dstDy := dst.Bounds().Dy()
			dstRatio := float64(dstDx) / float64(dstDy)
			var resized image.Image
			if srcRatio < dstRatio {
				resized = imaging.Resize(src, 0, dstDy, resizeFilter)
			} else {
				resized = imaging.Resize(src, dstDx, 0, resizeFilter)
			}
			draw.Draw(img, resized.Bounds(), resized, image.Point{}, draw.Src)
		})
		if err != nil {
			err = fmt.Errorf("draw xImg: %w", err)
			return
		}
	}
	err = c.Do(api.UserReq{Username: string(ann.Author)}, &ann.XAuthor)
	if err != nil {
		err = fmt.Errorf("author: %w", err)
		return
	}
	ann.XOrg, err = ann.ReqOrg(c)
	if err != nil {
		err = fmt.Errorf("org: %w", err)
		return
	}
	ann.XURL = ann.URL(c).String()

	img, err = drawIcon(c.HTTPClient(), img, ann.XOrg)
	if err != nil {
		return
	}
	_, err = drawText(img, image.Point{X: 128, Y: 128}, "abc")
	if err != nil {
		return nil, "", err
	}

	var text2 string
	text2, err = body(ann)
	if err != nil {
		return
	}
	text += text2
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

func drawText(orig draw.Image, point image.Point, text string) (image.Image, error) {
	ctx := freetype.NewContext()
	notoSans, err := truetype.Parse(font.NotoSansDisplayThingItalic)
	if err != nil {
		return nil, err
	}
	ctx.SetFont(notoSans)
	ctx.SetFontSize(10)
	ctx.SetSrc(&image.Uniform{C: image.White})
	ctx.SetDst(orig)
	_, err = ctx.DrawString(text, fixed.Point26_6{
		X: fixed.I(point.X),
		Y: fixed.I(point.Y),
	})
	if err != nil {
		return nil, err
	}
	return orig, nil
}

func drawIcon(hc api.HTTPClient, orig image.Image, org api.Org) (*image.NRGBA, error) {
	icon, _, err := getImage(hc, org.Icon)
	if err != nil {
		return nil, err
	}
	icon = imaging.Resize(
		icon,
		int(float64(orig.Bounds().Dx())*ratioOfIconToImage),
		0,
		resizeFilter,
	)
	mask := image.NewNRGBA(icon.Bounds())
	drawCircle(mask, image.Point{}, icon.Bounds().Dx(), color.Opaque)
	draw.DrawMask(icon, icon.Bounds(), icon, icon.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	return imaging.Paste(
		orig,
		icon,
		image.Point{
			Y: orig.Bounds().Max.Sub(icon.Bounds().Size()).Y,
		},
	), nil
}

func drawImage(img *image.NRGBA, hc api.HTTPClient, url string, draw func(dst, src image.Image)) (err error) {
	var annImg image.Image
	annImg, _, err = getImage(hc, url)
	if err != nil {
		err = fmt.Errorf("xImg: %w", err)
		return
	}
	draw(img, annImg)
	return nil
}

func getImage(hc api.HTTPClient, url string) (img draw.Image, name string, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := hc.Do(request)
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
