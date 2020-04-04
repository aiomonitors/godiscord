package godiscord

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

//Embed is a struct representing a Discord embed object
type Embed struct {
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Content   string         `json:"content"`
	Embeds    []EmbedElement `json:"embeds"`
}

//EmbedElement is a struct representing an Embed element of the Embed struct
type EmbedElement struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int64   `json:"color"`
	Fields      []Field `json:"fields"`
	Thumbnail   Image   `json:"thumbnail,omitempty"`
	Image       Image   `json:"image,omitempty"`
	Footer      Footer  `json:"footer"`
}

//Author represents the author of the embed
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

//Field represents a field in an embed
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

//Footer represents the footer of an embed
type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

//Image represents the image of an embed
type Image struct {
	URL string `json:"url"`
}

//Webhook represents a webhook
type Webhook struct {
	URL     string `json:"webhook"`
	IconURL string `json:"icon_url"`
	Text    string `json:"text"`
}

//NewEmbed creates a new embed object.
//Returns Embed.
func NewEmbed(Title, Description, URL string) Embed {
	e := Embed{}
	emb := EmbedElement{Title: Title, Description: Description, URL: URL}
	e.Embeds = append(e.Embeds, emb)
	return e
}

//SetAuthor sets the author of the Embed
func (e *Embed) SetAuthor(Name, URL, IconURL string) {
	if len(e.Embeds) == 0 {
		emb := EmbedElement{Author: Author{Name, URL, IconURL}}
		e.Embeds = append(e.Embeds, emb)
	} else {
		e.Embeds[0].Author = Author{Name, URL, IconURL}
	}
}

//SetColor takes in a hex code and sets the color of the Embed.
//Returns an error if the hex is invalid
func (e *Embed) SetColor(color string) error {
	color = strings.Replace(color, "0x", "", -1)
	color = strings.Replace(color, "0X", "", -1)
	color = strings.Replace(color, "#", "", -1)
	colorInt, err := strconv.ParseInt(color, 16, 64)
	if err != nil {
		return errors.New("Invalid hex code passed")
	}
	e.Embeds[0].Color = colorInt
	return nil
}

//SetThumbnail sets the thumbnail of the embed.
//Returns an error if the embed was not initialized properly
func (e *Embed) SetThumbnail(URL string) error {
	if len(e.Embeds) < 1 {
		return errors.New("Invalid Embed passed in, Embed.Embeds must have at least one EmbedElement")
	}
	e.Embeds[0].Thumbnail = Image{URL}
	return nil
}

//SetImage sets the image of the embed
//Returns an error if the embed was not initialized properly
func (e *Embed) SetImage(URL string) error {
	if len(e.Embeds) < 1 {
		return errors.New("Invalid Embed passed in, Embed.Embeds must have at least one EmbedElement")
	}
	e.Embeds[0].Image = Image{URL}
	return nil
}

//SetFooter sets the footer of the embed.
//Returns an error if the embed was not initialized properly
func (e *Embed) SetFooter(Text, IconURL string) error {
	if len(e.Embeds) < 1 {
		return errors.New("Invalid Embed passed in, Embed.Embeds must have at least one EmbedElement")
	}
	e.Embeds[0].Footer = Footer{Text, IconURL}
	return nil
}

//AddField adds a frield to the Embed.
//Returns an error if the embed was not initialized properly
func (e *Embed) AddField(Name, Value string, Inline bool) error {
	if len(e.Embeds) < 1 {
		return errors.New("Invalid Embed passed in, Embed.Embeds must have at least one EmbedElement")
	}
	e.Embeds[0].Fields = append(e.Embeds[0].Fields, Field{Name, Value, Inline})
	return nil
}

//SendToWebhook sents the Embed to a webhook.
//Returns error if embed was invalid or there was an error posting to the webhook.
func (e *Embed) SendToWebhook(Webhook string) error {
	embed, marshalErr := json.Marshal(e)
	if marshalErr != nil {
		return marshalErr
	}
	_, postErr := http.Post(Webhook, "application/json", bytes.NewBuffer(embed))
	if postErr != nil {
		return postErr
	}
	return nil
}
