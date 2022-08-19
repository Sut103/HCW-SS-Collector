package discord

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const base_url = "https://discord.com/api"
const api_resource = "/channels"

type ChannelAPI struct {
	BotToken        string
	ChannelEndpoint string
}

func NewChannelAPI() ChannelAPI {
	bot_token := os.Getenv("DISCORD_BOT_TOKEN")
	channel_id := os.Getenv("DISCORD_CHANNEL_ID")

	return ChannelAPI{
		BotToken:        strings.Join([]string{"Bot ", bot_token}, ""),
		ChannelEndpoint: strings.Join([]string{base_url, api_resource, "/", channel_id}, ""),
	}
}

type ChannelMessage struct {
	Attachments []struct {
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
		Height      int    `json:"height"`
		ID          string `json:"id"`
		ProxyURL    string `json:"proxy_url"`
		Size        int    `json:"size"`
		URL         string `json:"url"`
		Width       int    `json:"width"`
	} `json:"attachments"`
	Author struct {
		Avatar           string      `json:"avatar"`
		AvatarDecoration interface{} `json:"avatar_decoration"`
		Discriminator    string      `json:"discriminator"`
		ID               string      `json:"id"`
		PublicFlags      int         `json:"public_flags"`
		Username         string      `json:"username"`
	} `json:"author"`
	ChannelID       string        `json:"channel_id"`
	Components      []interface{} `json:"components"`
	Content         string        `json:"content"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Embeds          []interface{} `json:"embeds"`
	Flags           int           `json:"flags"`
	ID              string        `json:"id"`
	MentionEveryone bool          `json:"mention_everyone"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Mentions        []interface{} `json:"mentions"`
	Pinned          bool          `json:"pinned"`
	Timestamp       time.Time     `json:"timestamp"`
	Tts             bool          `json:"tts"`
	Type            int           `json:"type"`
}

func (ca *ChannelAPI) getChannelMessage(query map[string]string) ([]ChannelMessage, error) {
	req, err := http.NewRequest("GET", strings.Join([]string{ca.ChannelEndpoint, "/messages"}, ""), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", ca.BotToken)

	params := req.URL.Query()
	for key, value := range query {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var channelMessages []ChannelMessage
	err = json.Unmarshal(bytes, &channelMessages)
	if err != nil {
		return nil, err
	}

	return channelMessages, nil
}

func (ca *ChannelAPI) GetChannelMessageAll() ([]ChannelMessage, error) {
	cm, err := ca.getChannelMessage(
		map[string]string{
			"limit": "100",
		})
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func (ca *ChannelAPI) GetChannelMessagesNewer(message_id string) ([]ChannelMessage, error) {
	cm, err := ca.getChannelMessage(
		map[string]string{
			"after": message_id,
		})
	if err != nil {
		return nil, err
	}

	return cm, nil
}
