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

func TestResolverSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}
