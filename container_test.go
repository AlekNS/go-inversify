package inversify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ContainerTestSuite struct {
	suite.Suite
}

func (t *ContainerTestSuite) TestBasic() {
	c1 := Container()
	c1.Bind(1).To("val1")
	c1.Build()

	c2 := Container()
	c2.Bind(2).To("val2")
	c2.Build()

	c3 := c1.Merge(c2)
	t.True(c3.IsBound(1))
	t.True(c3.IsBound(2))
	t.False(c3.IsBound(3))
}

func (t *ContainerTestSuite) TestGet() {
	c1 := Container()
	c1.Bind(1).To("V1")
	c1.Build()

	c2 := Container()
	c2.Bind(2).ToTypedFactory(func(i1 string) (string, error) {
		return fmt.Sprintf("V2(1:%s)", i1), nil
	}, 1)
	c2.Bind(3).ToFactory(func(i1 Any, i2 Any, i3 Any) (Any, error) {
		return fmt.Sprintf("V3(1:%s,2:%s,3:%v)", i1, i2, i3), nil
	}, 1, 2, Optional(1000))
	c2.SetParent(c1)
	c2.Build()

	val, err := c2.Get(3)
	t.NoError(err)
	t.Equal("V3(1:V1,2:V2(1:V1),3:<nil>)", val)
}

func TestContainerSute(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
