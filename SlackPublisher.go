package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Text        string        `json:"text"`
	Channel     string        `json:"channel,omitempty"`
	UserName    string        `json:"username,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback,omitempty"` // plain text summary
	Color      string  `json:"color,omitempty"`    // {good|warning|danger|hex}
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Title      string  `json:"title,omitempty"` // larger, bold text at top of attachment
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	FooterIcon string  `json:"footer,omitempty"`
	Footer     string  `json:"footer_icon,omitempty"`
	Timestamp  int     `json:"ts,omitempty"` // Unix timestamp
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

type Client struct {
	url        string
	HTTPClient Poster
}

// New Slack Incoming WebHook Client using http.DefaultClient for its Poster.
func New(url string) *Client {
	return &Client{url: url, HTTPClient: http.DefaultClient}
}

type Poster interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

func (c *Client) SimpleToChannel(msg string, channel string) error {
	return c.Send(&Message{Text: msg, Channel: channel})
}

// Simple text message.
func (c *Client) Simple(msg string) error {
	return c.Send(&Message{Text: msg})
}

func (c *Client) Send(msg *Message) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Post(c.url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard response body to reuse connection
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
