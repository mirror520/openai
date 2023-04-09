package chat

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type chatTestSuite struct {
	suite.Suite
	c *Chat
}

func (suite *chatTestSuite) SetupSuite() {
	opts := new(Options)
	opts.Temperature = new(float64)
	*opts.Temperature = 0.0

	suite.c = NewChat("gpt-3.5-turbo", "You are a helpful assistant.", opts)
}

func (suite *chatTestSuite) TestRequest() {
	suite.c.AddMessage(&Message{
		Role:    User,
		Content: "Hello!",
	})
	req := suite.c.Request()

	suite.Equal("gpt-3.5-turbo", req.Model)
	suite.Len(req.Messages, 2)

	suite.Equal(0.0, *req.Options.Temperature)

	req.Stream = new(bool)
	*req.Stream = true

	suite.Nil(suite.c.Options.Stream)
	suite.True(*req.Options.Stream)
}

func (suite *chatTestSuite) TestUpdateOptions() {
	opts := new(Options)
	opts.MaxTokens = new(int)
	*opts.MaxTokens = 10

	err := suite.c.Options.Update(opts)
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	suite.Equal(10, *suite.c.MaxTokens)
	suite.Equal(0.0, *suite.c.Temperature)
}

func TestChatTestSuite(t *testing.T) {
	suite.Run(t, new(chatTestSuite))
}
