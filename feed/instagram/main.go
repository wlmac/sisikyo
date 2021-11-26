package instagram

import (
	"bytes"
	"fmt"
	"image/png"
	"net/url"
	"sync"
	"time"

	"github.com/ahmdrz/goinsta/v2"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/feed"
)

// Instagram fulfills feed.Sink.
type Instagram struct {
	c      *api.Client
	fmt    feed.FormatImager
	insta  goinsta.Instagram
	delay  func() time.Duration
	noCopy [0]sync.Mutex
}

var _ feed.Sink = (*Instagram)(nil)

func (i *Instagram) wait() {
	time.Sleep(i.delay())
}

// Post fulfills feed.Sink.
func (i *Instagram) Post(ann api.Ann) (url *url.URL, err error) {
	img, caption, err := i.fmt.FormatImage(ann)
	if err != nil {
		return nil, fmt.Errorf("fmt: %w", err)
	}
	imgBuf := new(bytes.Buffer)
	err = png.Encode(imgBuf, img)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	i.wait()
	err = i.insta.Login()
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	i.wait()
	_, err = i.insta.UploadPhoto(imgBuf, caption, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("upload: %w", err)
	}
	return nil, nil
}
