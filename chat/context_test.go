package chat

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type contextTestSuite struct {
	suite.Suite
	ctx *Context
}

func (suite *contextTestSuite) SetupSuite() {
	opts := new(Options)
	opts.Temperature = new(float64)
	*opts.Temperature = 0.0

	suite.ctx = NewContext("gpt-3.5-turbo", "You are a helpful assistant.", opts)
}

func (suite *contextTestSuite) TestRequest() {
	suite.ctx.AddMessage(&Message{
		Role:    User,
		Content: "Hello!",
	})
	req := suite.ctx.Request()

	suite.Equal("gpt-3.5-turbo", req.Model)
	suite.Len(req.Messages, 2)

	suite.Equal(0.0, *req.Options.Temperature)

	req.Stream = new(bool)
	*req.Stream = true

	suite.Nil(suite.ctx.Options.Stream)
	suite.True(*req.Options.Stream)
}

func (suite *contextTestSuite) TestUpdateOptions() {
	opts := new(Options)
	opts.MaxTokens = new(int)
	*opts.MaxTokens = 10

	err := suite.ctx.Options.Update(opts)
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	suite.Equal(10, *suite.ctx.MaxTokens)
	suite.Equal(0.0, *suite.ctx.Temperature)
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(contextTestSuite))
}
