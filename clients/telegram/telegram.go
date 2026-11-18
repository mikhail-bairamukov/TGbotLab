package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "GetUpdates"
	sendMessageMethod = "SendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// TODO: ? type baseBotPath string ?
func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) GetUpdates(offset int, limit int) ([]Update, error) {
	const op = "GetUpdates"

	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, "c.doRequest", err)
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, "json.Unmarshal", err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	const op = "SendMessage"

	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, query)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", op, "c.doRequest", err)
	}

	return nil
}

// TODO: response status "ok" is not handled
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const op = "doRequest"

	// u -> requestURL
	requestURL := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, "http.NewRequest", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, "c.client.Do", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, "read response body", err)
	}

	return body, nil
}
