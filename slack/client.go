package slack

import (
	"fmt"

	slackClient "github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client is a slack client
type Client struct {
	slack  *slackClient.Client
	logger *logrus.Logger
}

// New returns a webhook client
func New(token string, logger *logrus.Logger) *Client {
	client := slackClient.New(token)
	return &Client{
		slack:  client,
		logger: logger,
	}
}

//GetSlackChannelID returns the chanel id from an email
func (c *Client) GetSlackChannelID(email string) (string, error) {
	user, err := c.slack.GetUserByEmail(email)
	if err != nil {
		return "", errors.Wrap(err, "could not find slack user for email")
	}
	if user == nil {
		return "", errors.New("email not found")
	}
	c.logger.Info(fmt.Sprintf("userID: %s", user.ID))
	_, _, channelID, err := c.slack.OpenIMChannel(user.ID)
	return channelID, errors.Wrap(err, "could not open dm channel with user")
}

// PostMessage posts a message
func (c *Client) PostMessage(message Message) error {
	channelID, err := c.GetSlackChannelID(message.Email)
	if err != nil {
		return err
	}
	params := slackClient.PostMessageParameters{
		UnfurlLinks: true,
		Attachments: message.Attachments,
	}
	_, _, err = c.slack.PostMessage(channelID, message.Text, params)
	return errors.Wrap(err, "could not post message")
}