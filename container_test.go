package inversify

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testInterface1 interface {
}

type testInterface2 interface {
}

type ContainerTestSuite struct {
	suite.Suite
}

const (
	testDep1 = iota
	testDep2
	testDep3
	testOtherDep
)

func (t *ContainerTestSuite) TestBasic() {
	c1 := NewContainer()
	c1.Bind(testDep1).To(1000)
	c1.Build()

	t.True(c1.IsBound(testDep1))
	value, _ := c1.Get(testDep1)
	t.Equal(1000, value)

	c1.Unbind(testDep1)
	c1.Build()
	t.False(c1.IsBound(testDep1))
}

func (t *ContainerTestSuite) TestFallthroughError() {
	c1 := NewContainer()
	c1.Bind((*testInterface1)(nil)).ToFactory(func() (Any, error) {
		return nil, errors.New("FactoryError")
	})
	c1.Bind(testDep2).ToTypedFactory(func(t1 testInterface1) (string, error) {
		return "Not working", nil
	}, (*testInterface1)(nil))
	c1.Build()

	t.True(c1.IsBound(testDep2))
	value, err := c1.Get(testDep2)
	t.EqualError(err, "FactoryError")
	t.Nil(value)
}

func (t *ContainerTestSuite) TestMerge() {
	c1 := NewContainer()
	c1.Bind(testDep1).To("val1")
	c1.Bind((*testInterface1)(nil)).To("val3")
	c1.Bind((*testInterface2)(nil)).To("val4")
	c1.Bind(testOtherDep).ToTypedFactory(func(t1, t2 string) (string, error) {
		return t1 + t2, nil
	}, (*testInterface1)(nil), (*testInterface2)(nil))
	c1.Build()

	v, _ := c1.Get((*testInterface2)(nil))
	t.Equal("val4", v)

	c2 := NewContainer()
	c2.Bind(testDep2).To("val2")
	c2.Build()

	c3 := c1.Merge(c2)
	c3.Build()
	t.True(c3.IsBound(testDep1))
	t.True(c3.IsBound(testDep2))
	t.False(c3.IsBound(testDep3))

	t.True(c3.IsBound(testOtherDep))

	value, _ := c3.Get(testOtherDep)
	t.Equal("val3val4", value)
}

func (t *ContainerTestSuite) TestParent() {
	c1 := NewContainer()

	c1.Bind(testDep1).To("V1")
	c1.Bind(testDep2).ToTypedFactory(func(i1 string, any Any) (string, error) {
		t.Nil(any)
		return fmt.Sprintf("V2(1:%s)", i1), nil
	}, testDep1, Optional(testOtherDep))

	c1.Build()

	c2 := NewContainer()
	c2.Bind(testDep3).ToFactory(func(i1 Any, i2 Any, i3 Any) (Any, error) {
		t.Nil(i3)
		return fmt.Sprintf("V3(1:%s,2:%s,3:%v)", i1, i2, i3), nil
	}, testDep1, testDep2, Optional(testOtherDep))

	t.Nil(c2.GetParent())
	c2.SetParent(c1)
	t.Equal(c1, c2.GetParent())

	c2.Build()

	val, err := c2.Get(testDep3)
	t.NoError(err)
	t.Equal("V3(1:V1,2:V2(1:V1),3:<nil>)", val)
}

func TestContainerSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
