package discord

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/feed"
)

// Discord fulfills feed.Sink.
type Discord struct {
	c      *api.Client
	s      *discordgo.Session
	noCopy [0]sync.Mutex
}

// TODO: render text for insta testing

var _ feed.Sink = (*Discord)(nil)

// NewDiscord returns a new Discord.
func NewDiscord(token string, c *api.Client) (*Discord, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Discord{c: c, s: s}, nil
}

// Close releases resources.
func (d *Discord) Close() error { return d.s.Close() }

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
func (d *Discord) Post(ann api.Ann) (u *url.URL, err error) {
	var org api.Org
	orgs := api.OrgsResp{}
	err = d.c.Do(api.OrgsReq{}, &orgs)
	if err != nil {
		return nil, fmt.Errorf("orgs: %w", err)
	}
	for _, org2 := range orgs {
		if org2.Name == ann.Org {
			org = org2
		}
	}

	author := api.UserResp{}
	err = d.c.Do(api.UserReq{Username: string(ann.Author)}, &author)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	xImageURL, err := url.Parse(ann.XImageURL)
	if err != nil {
		return nil, fmt.Errorf("xImg parse: %w", err)
	}
	orgBanner, err := url.Parse(ann.XOrg.Banner)
	if err != nil {
		return nil, fmt.Errorf("org banner parse: %w", err)
	}

	msg := &discordgo.MessageEmbed{
		URL:         d.c.Rel(ann.URL(d.c)).String(),
		Type:        discordgo.EmbedTypeArticle,
		Title:       ann.Title,
		Description: ann.Body,
		Timestamp:   ann.LastModified.Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("By %s %s (%s)", author.FirstName, author.LastName, author.Username),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: d.c.Rel(xImageURL).String(),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: d.c.Rel(orgBanner).String(),
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:     d.c.Rel(org.URL(d.c)).String(),
			Name:    org.Name,
			IconURL: org.Icon,
		},
		Fields: []*discordgo.MessageEmbedField{
			tagsField(ann.Tags),
			{Name: "Audience", Value: map[bool]string{true: "everyone", false: "members"}[ann.Public], Inline: true},
		},
	}
	sent, err := d.s.ChannelMessageSendEmbed("892881405862871070", msg)
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
