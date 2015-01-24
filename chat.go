package slack

import (
	"errors"
	"net/url"
	"strings"
)

const (
	DEFAULT_MESSAGE_USERNAME     = ""
	DEFAULT_MESSAGE_PARSE        = ""
	DEFAULT_MESSAGE_LINK_NAMES   = 0
	DEFAULT_MESSAGE_UNFURL_LINKS = true
	DEFAULT_MESSAGE_UNFURL_MEDIA = false
	DEFAULT_MESSAGE_ICON_URL     = ""
	DEFAULT_MESSAGE_ICON_EMOJI   = ""
)

type chatResponseFull struct {
	ChannelID string `json:"channel"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackResponse
}

type PostMessageParameters struct {
	Text        string
	Username    string
	Parse       string
	LinkNames   int
	Attachments []byte
	UnfurlLinks bool
	UnfurlMedia bool
	IconURL     string
	IconEmoji   string
}

func NewPostMessageParameters() PostMessageParameters {
	return PostMessageParameters{
		Username:    DEFAULT_MESSAGE_USERNAME,
		Parse:       DEFAULT_MESSAGE_PARSE,
		LinkNames:   DEFAULT_MESSAGE_LINK_NAMES,
		Attachments: nil,
		UnfurlLinks: DEFAULT_MESSAGE_UNFURL_LINKS,
		UnfurlMedia: DEFAULT_MESSAGE_UNFURL_MEDIA,
		IconURL:     DEFAULT_MESSAGE_ICON_URL,
		IconEmoji:   DEFAULT_MESSAGE_ICON_EMOJI,
	}
}

func chatRequest(path string, values url.Values, debug bool) (*chatResponseFull, error) {
	response := &chatResponseFull{}
	err := ParseResponse(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func (api *Slack) DeleteMessage(channelId, messageTimestamp string) (string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"ts":      {messageTimestamp},
	}
	response, err := chatRequest("chat.delete", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.ChannelID, response.Timestamp, nil
}

func EscapeMessage(message string) string {
	/*
		& replaced with &amp;
		< replaced with &lt;
		> replaced with &gt;
	*/
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(message)
}

func (api *Slack) PostMessage(channelId string, text string, params PostMessageParameters) (string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"text":    {EscapeMessage(text)},
	}
	if params.Username != DEFAULT_MESSAGE_USERNAME {
		values.Set("username", string(params.Username))
	}
	if params.Parse != DEFAULT_MESSAGE_PARSE {
		values.Set("parse", string(params.Parse))
	}
	if params.LinkNames != DEFAULT_MESSAGE_LINK_NAMES {
		values.Set("link_names", "1")
	}
	if params.Attachments != nil {
		values.Set("attachments", string(params.Attachments))
	}
	if params.UnfurlLinks == DEFAULT_MESSAGE_UNFURL_LINKS {
		values.Set("unfurl_links", "false")
	}
	if params.UnfurlMedia != DEFAULT_MESSAGE_UNFURL_MEDIA {
		values.Set("unfurl_media", "true")
	}
	if params.IconURL != DEFAULT_MESSAGE_ICON_URL {
		values.Set("icon_url", params.IconURL)
	}
	if params.IconEmoji != DEFAULT_MESSAGE_ICON_EMOJI {
		values.Set("icon_emoji", params.IconEmoji)
	}

	response, err := chatRequest("chat.postMessage", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.ChannelID, response.Timestamp, nil
}

func (api *Slack) UpdateMessage(channelId, timestamp, text string) (string, string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channelId},
		"text":    {EscapeMessage(text)},
		"ts":      {timestamp},
	}
	response, err := chatRequest("chat.update", values, api.debug)
	if err != nil {
		return "", "", "", err
	}
	return response.ChannelID, response.Timestamp, response.Text, nil

}