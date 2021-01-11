package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	defaultUser = "doorbot"
	defaultIcon = ":robot:"
)

type Slack struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Channel   string `json:"channel"`
}

func SendMsg(slackURL, msg string) error {
	params := Slack{
		Username:  defaultUser,
		IconEmoji: defaultIcon,
		Text:      msg,
	}
	return Send(slackURL, params)
}

func Send(slackURL string, params Slack) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	resp, err := http.PostForm(
		slackURL,
		url.Values{"payload": {string(data)}},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
