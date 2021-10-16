package text

import (
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
	"image"
)

type DrawTexter interface {
	DrawText(img image.Image, text string) error
}

type FreeType struct {
	ctx freetype.Context
}

func (f *FreeType) DrawText(img image.Image, point fixed.Point26_6, text string) error {
	f.ctx.SetSrc(img)
	_, err := f.ctx.DrawString(text, point)
	if err != nil {
		return err
	}
	return nil
}
