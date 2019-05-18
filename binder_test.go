package inversify

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

const resolvedValue = "resolved value"

type BindingTestSuite struct {
	suite.Suite

	binding *Binding
}

func defaultDepVal(val Any) FactoryFunc     { return func() (Any, error) { return val, nil } }
func defaultDepError(err error) FactoryFunc { return func() (Any, error) { return nil, err } }

func (t *BindingTestSuite) BeforeTest(suiteName, testName string) {
	t.binding = &Binding{sync.Once{}, nil, NAny{}, NAny{}}
	t.binding.resolves = NAny{defaultDepVal("resolved"), defaultDepVal(" value")}
}

func (t *BindingTestSuite) TestBindValue() {
	binding := t.binding

	binding.To(resolvedValue)
	val, err := binding.factory()

	t.Equal(resolvedValue, val)
	t.NoError(err)
}

func (t *BindingTestSuite) TestBindAbstractFactory() {
	binding := t.binding

	binding.ToFactory(func(a, b Any) (Any, error) {
		return a.(string) + b.(string), nil
	}, testDep1, testDep2)

	val, err := binding.factory()

	t.Equal(resolvedValue, val)
	t.NoError(err)
}

func (t *BindingTestSuite) TestBindAbstractFactoryError() {
	binding := t.binding

	binding.ToFactory(func(a, b Any) (Any, error) {
		return a.(string) + b.(string), nil
	}, testDep1, testDep2)
	binding.resolves = NAny{defaultDepVal("resolved"), defaultDepError(errors.New("error"))}
	val, err := binding.factory()

	t.Nil(val)
	t.Error(err)
}

func (t *BindingTestSuite) TestBindTypedFactorySimple() {
	binding := t.binding

	counter := 0
	binding.ToTypedFactory(func(a, b string) (string, error) {
		counter++
		return a + b, nil
	}, 1, 2)
	val, err := binding.factory()
	val, err = binding.factory()

	t.Equal(resolvedValue, val)
	t.NoError(err)
	t.Equal(2, counter)
}

type interfaceValue interface {
	get() string
}
type structValue1 struct{}

func (v structValue1) get() string {
	return "resolved"
}

func (t *BindingTestSuite) TestBindTypedFactoryAndPointers() {
	binding := t.binding

	binding.ToTypedFactory(func(a interfaceValue, b *string) (string, error) {
		return a.get() + *b, nil
	}, testDep1, testDep2)
	valueString := " value"
	binding.resolves = NAny{defaultDepVal(&structValue1{}), defaultDepVal(&valueString)}
	val, err := binding.factory()

	t.Equal(resolvedValue, val)
	t.NoError(err)
}

func (t *BindingTestSuite) TestBindTypedFactorySingleton() {
	binding := t.binding

	counter := 0
	binding.ToTypedFactory(func(a, b string) (string, error) {
		counter++
		return a + b, nil
	}, testDep1, testDep2).InSingletonScope()
	val, err := binding.factory()
	val, err = binding.factory()

	t.Equal(resolvedValue, val)
	t.NoError(err)
	t.Equal(1, counter)
}

func (t *BindingTestSuite) TestBindTypedFactoryNoDeps() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	binding.ToTypedFactory(func() (string, error) {
		return "no dependencies", nil
	})
	val, err := binding.factory()

	t.Equal("no dependencies", val)
	t.NoError(err)
}

func TestBindingSuite(t *testing.T) {
	suite.Run(t, new(BindingTestSuite))
}
