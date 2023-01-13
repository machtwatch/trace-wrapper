package tracer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AttributeTestSuite struct {
	suite.Suite
}

func TestAttributeTestSuite(t *testing.T) {
	suite.Run(t, new(AttributeTestSuite))
}

func (a *AttributeTestSuite) TestSetAttribute() {
	attr := SetAttribute("key", "value")
	a.True(attr.getAttribute().Valid())
}
