package debug

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/feed"
	"gitlab.com/mirukakoro/sisikyo/feed/image2"
)

// Debug is a feed.Sink just for debugging purposes.
type Debug struct {
	c      *api.Client
	s      *discordgo.Session
	noCopy [0]sync.Mutex
}

// TODO: render text for insta testing

var _ feed.Sink = (*Debug)(nil)

// NewDebug makes a new Debug.
func NewDebug(token string, c *api.Client) (*Debug, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Debug{c: c, s: s}, nil
}

// Close releases resources.
func (d *Debug) Close() error { return d.s.Close() }

func tagsField(tags []api.Tag) *discordgo.MessageEmbedField {
	value := ""
	for i, tag := range tags {
		if i != 0 {
			value += ", "
		}
		value += tag.Name
	}
	return &discordgo.MessageEmbedField{
		Name:  "Tags",
		Value: value,
	}
}

const timeFmt = `2006-01-02T15:04:05`

// Post fulfills feed.Sink.
func (d *Debug) Post(ann api.Ann) (u *url.URL, err error) {
	svg, caption, err := image2.AnnToImageAndText2(d.c, ann)
	if err != nil {
		return nil, fmt.Errorf("img: %w", err)
	}

	msg := &discordgo.MessageSend{
		Content: fmt.Sprint(caption),
		Files: []*discordgo.File{
			{
				Name:        "image.svg",
				ContentType: "image/svg+xml",
				Reader:      bytes.NewBufferString(svg),
			},
		},
	}

	sent, err := d.s.ChannelMessageSendComplex("892881405862871070", msg)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	return u.Parse(fmt.Sprintf(
		"https://discord.com/channels/%s/%s/%s",
		sent.GuildID,
		sent.ChannelID,
		sent.ID,
	))
}
