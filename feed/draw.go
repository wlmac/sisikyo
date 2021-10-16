package feed

import (
	"image"
	"image/color"
	"image/draw"
)

// drawCircle draws a circle on img starting from sp and with radius r.
// code used in function is by icza <https://stackoverflow.com/users/1705598/icza> from StackOverflow <https://stackoverflow.com/a/51627141>
func drawCircle(img draw.Image, sp image.Point, r int, c color.Color) {
	x0, y0 := sp.X, sp.Y
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
