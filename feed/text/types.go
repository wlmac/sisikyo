package text

import (
	"image"

	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
)

// DrawTexter draws text onto an Image.
type DrawTexter interface {
	DrawText(img image.Image, text string) error
}

// FreeType wraps freetype.Context.
type FreeType struct {
	ctx freetype.Context
}

// DrawText fulfills DrawTexter.
func (f *FreeType) DrawText(img image.Image, point fixed.Point26_6, text string) error {
	f.ctx.SetSrc(img)
	_, err := f.ctx.DrawString(text, point)
	if err != nil {
		return err
	}
	return nil
}
