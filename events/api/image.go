package api

import (
	"regexp"
	"strings"
)

// TODO: get an image URL from an announcement in this order:
//     - first image in body
//     - pfp and banner of posting org

var mdImagePattern, _ = regexp.Compile(`\!\[[^\]]*\]\([^\)]*\)`)

func GetImageFromMd(src string) (alt, url string, found bool) {
	mdImage := mdImagePattern.FindString(src)
	if mdImage != "" {
		found = true
		split := strings.SplitN(mdImage, "](", 2)
		alt = split[0][2:]
		url = split[1][:len(split[1])-1]
	}
	return
}
