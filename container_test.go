package inversify

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type structure struct {
	val int
}

type testInterface1 interface {
	get1() int
}

type testInterface1Impl struct{}

func (impl testInterface1Impl) get1() int { return 1 }

type testInterface2 interface {
	get2() int
}

type testInterface2Impl struct{}

func (impl testInterface2Impl) get2() int { return 2 }

type ContainerTestSuite struct {
	suite.Suite
}

const (
	testDep1 = iota
	testDep2
	testDep3
	testOtherDep1
	testOtherDep2
)

func (t *ContainerTestSuite) TestBasic() {
	c1 := NewContainer("basic")
	c1.Bind(testDep1).To(resolvedValue)
	c1.Build()

	t.True(c1.IsBound(testDep1))

	value, _ := c1.Get(testDep1)

	t.Equal(resolvedValue, value)

	c1.Unbind(testDep1)
	c1.Build()

	t.False(c1.IsBound(testDep1))
}

type depType1 int
type depType2 int

func (t *ContainerTestSuite) TestBasicTypes() {
	var dep1 depType1 = 1
	var dep2 depType2 = 1

	c1 := NewContainer("")
	c1.Bind(dep1).To("val1")
	c1.Bind(dep2).To("val2")
	c1.Build()

	val1, _ := c1.Get(dep1)
	val2 := c1.MustGet(dep2)

	t.Equal("val1", val1)
	t.Equal("val2", val2)
}

func (t *ContainerTestSuite) TestNamedBasic() {
	c1 := NewContainer("withNames")
	c1.Bind(testDep1).To(resolvedValue)
	c1.Bind(testDep1, "other").To(resolvedValue + resolvedValue)
	c1.Build()

	t.True(c1.IsBound(testDep1))

	value, _ := c1.Get(testDep1)
	valueOther, _ := c1.Get(testDep1, "other")

	t.Equal(resolvedValue, value)
	t.Equal(resolvedValue+resolvedValue, valueOther)

	c1.Unbind(testDep1)
	c1.Build()

	t.False(c1.IsBound(testDep1))
}

func (t *ContainerTestSuite) TestNamedCycle() {
	c1 := NewContainer("withNamesAndCycle")
	c1.Bind(testDep1).ToFactory(func(i Any) (Any, error) {
		return 1, nil
	}, Named(testDep3, "named"))
	c1.Bind(testDep2).ToFactory(func(i Any) (Any, error) {
		return 2, nil
	}, testDep1)
	c1.Bind(testDep3, "named").ToFactory(func(i Any) (Any, error) {
		return 3, nil
	}, testDep2)

	t.Panics(func() { // the container has cycle dependencies: [0[] 1[] 2[named]]
		c1.Build()
	})
}

func (t *ContainerTestSuite) TestFallthroughError() {
	c1 := NewContainer("withError")
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
	c1 := NewContainer("src1")
	c1.Bind(testDep1).To("val1")
	c1.Bind((*testInterface1)(nil)).To(&testInterface1Impl{})
	c1.Bind((*testInterface2)(nil)).To(&testInterface2Impl{})
	c1.Bind((*structure)(nil)).To(&structure{})
	c1.Bind(testOtherDep1).ToFactory(func(t1 Any) (Any, error) {
		t.Equal(0, t1.(*structure).val)
		return "val3val4", nil
	}, (*structure)(nil))
	c1.Bind(testOtherDep2).ToTypedFactory(func(t1 *structure, t2 testInterface1) (string, error) {
		t.Equal(0, t1.val)
		t.Equal(1, t2.get1())
		return "val3val4", nil
	}, (*structure)(nil), (*testInterface1)(nil))

	c2 := NewContainer("src2")
	c2.Bind(testDep2).To("val2")

	c3 := c1.Merge(c2, "dst")
	c3.Build()

	t.True(c3.IsBound(testDep1))
	t.True(c3.IsBound(testDep2))
	t.False(c3.IsBound(testDep3))

	t.True(c3.IsBound(testOtherDep1))
	t.True(c3.IsBound(testOtherDep2))

	value, _ := c3.Get(testOtherDep1)

	t.Equal("val3val4", value)
}

func (t *ContainerTestSuite) TestParent() {
	c1 := NewContainer("parent")

	c1.Bind(testDep1).To("V1")
	c1.Bind(testDep2).ToTypedFactory(func(i1 string, any Any) (string, error) {
		t.Nil(any)
		return fmt.Sprintf("V2(1:%s)", i1), nil
	}, testDep1, Optional(Named(testOtherDep2, "test")))

	c2 := NewContainer("child")
	c2.Bind(testDep3).ToFactory(func(i1 Any, i2 Any, i3 Any) (Any, error) {
		t.Nil(i3)
		return fmt.Sprintf("V3(1:%s,2:%s,3:%v)", i1, i2, i3), nil
	}, testDep1, testDep2, Optional(testOtherDep1))

	t.Nil(c2.GetParent())

	c2.SetParent(c1)

	t.Equal(c1, c2.GetParent())

	c2.Build()

	val, err := c2.Get(testDep2)
	t.NoError(err)
	t.Equal("V2(1:V1)", val)

	val, err = c2.Get(testDep3)

	t.NoError(err)
	t.Equal("V3(1:V1,2:V2(1:V1),3:<nil>)", val)
}

func TestContainerSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}

func BenchmarkContainerGet(b *testing.B) {
	c1 := NewContainer("")
	c1.Bind(testDep1).To("val1")
	c1.Build()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c1.Get(testDep1)
	}
	b.StopTimer()

	b.ReportAllocs()
}

func BenchmarkContainerGetHierarchy(b *testing.B) {
	c1 := NewContainer("")

	c1.Bind(testDep1).To("V1")
	c1.Bind(testDep2).To("V2")

	c1.Build()

	c2 := NewContainer("")
	c2.Bind(testDep3).ToFactory(func(i1 Any, i2 Any, i3 Any) (Any, error) {
		return i1, nil
	}, testDep1, testDep2, Optional(testOtherDep1))

	c2.SetParent(c1)

	c2.Build()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c2.Get(testDep3)
	}
	b.StopTimer()

	b.ReportAllocs()
}
