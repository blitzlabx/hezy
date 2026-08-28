package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	system     string
}

func New(systemPrompt string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		system:     systemPrompt,
	}
}

type chatResponse struct {
	Status   bool   `json:"status"`
	Response string `json:"response"`
	State    string `json:"state"`
}

func (c *Client) Chat(prompt string, history []string) (string, error) {
	full := c.system + "\n\n"
	for _, h := range history {
		full += h + "\n"
	}
	full += "User: " + prompt

	u := "https://prexzyapis.com/ai/ch?q=" + url.QueryEscape(full)
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return string(body), nil
	}
	if cr.Response != "" {
		return cr.Response, nil
	}
	return strings.TrimSpace(string(body)), nil
}

func (c *Client) ChatWithState(prompt, state string) (string, string, error) {
	u := "https://prexzyapis.com/ai/askgpt5?prompt=" + url.QueryEscape(prompt)
	if state != "" {
		u += "&state=" + url.QueryEscape(state)
	}
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return string(body), "", nil
	}
	return cr.Response, cr.State, nil
}

type imageResponse struct {
	Status   bool `json:"status"`
	ImageURL []struct {
		Image struct {
			URL string `json:"url"`
		} `json:"image"`
	} `json:"image_url"`
}

func (c *Client) GenerateImage(prompt string) (string, error) {
	u := "https://prexzyapis.com/ai/dalle?prompt=" + url.QueryEscape(prompt)
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ir imageResponse
	if err := json.Unmarshal(body, &ir); err != nil {
		return "", fmt.Errorf("parse image response: %w", err)
	}
	if len(ir.ImageURL) > 0 && ir.ImageURL[0].Image.URL != "" {
		return ir.ImageURL[0].Image.URL, nil
	}
	return "", fmt.Errorf("no image returned")
}

func (c *Client) FootballScores() (string, error) {
	resp, err := c.httpClient.Get("https://prexzyapis.com/sports/football")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
