package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResolverTestSuite struct {
	suite.Suite
}

func (t *ResolverTestSuite) TestBasic() {
}

func TestResolverSute(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}
