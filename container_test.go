package inversify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type test1 interface {
}

type test2 interface {
}

type ContainerTestSuite struct {
	suite.Suite
}

func (t *ContainerTestSuite) TestBasic() {
	c1 := NewContainer()
	c1.Bind(1).To("val1")
	c1.Bind((*test1)(nil)).To("val3")
	c1.Bind((*test2)(nil)).To("val4")
	c1.Bind(10).ToTypedFactory(func(t1, t2 string) (string, error) {
		return t1 + t2, nil
	}, (*test1)(nil), (*test2)(nil))

	v, _ := c1.Get((*test2)(nil))
	t.Equal("val4", v)

	c2 := NewContainer()
	c2.Bind(2).To("val2")

	c3 := c1.Merge(c2)
	t.True(c3.IsBound(1))
	t.True(c3.IsBound(2))
	t.False(c3.IsBound(3))
}

func (t *ContainerTestSuite) TestGet() {
	c1 := NewContainer()

	c1.Bind(1).To("V1")

	c2 := NewContainer()
	c2.Bind(2).ToTypedFactory(func(i1 string) (string, error) {
		return fmt.Sprintf("V2(1:%s)", i1), nil
	}, 1)
	c2.Bind(3).ToFactory(func(i1 Any, i2 Any, i3 Any) (Any, error) {
		return fmt.Sprintf("V3(1:%s,2:%s,3:%v)", i1, i2, i3), nil
	}, 1, 2, Optional(1000))
	c2.SetParent(c1)

	val, err := c2.Get(3)
	t.NoError(err)
	t.Equal("V3(1:V1,2:V2(1:V1),3:<nil>)", val)
}

func TestContainerSute(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
