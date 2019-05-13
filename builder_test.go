package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactoryTestSuite struct {
	suite.Suite
}

func (t *FactoryTestSuite) TestBasic() {
}

func TestFactorySute(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}
