package inversify

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BindingTestSuite struct {
	suite.Suite
}

func defaultDepVal(val Any) FactoryFunc     { return func() (Any, error) { return val, nil } }
func defaultDepError(err error) FactoryFunc { return func() (Any, error) { return nil, err } }

func (t *BindingTestSuite) TestBindValue() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	binding.To("hello")
	val, err := binding.factory()

	t.Equal("hello", val)
	t.NoError(err)
}

func (t *BindingTestSuite) TestBindAbstractFactory() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	binding.ToFactory(func(a, b Any) (Any, error) {
		return a.(string) + b.(string), nil
	}, 1, 2)
	binding.resolves = NAny{defaultDepVal("1"), defaultDepVal("2")}
	val, err := binding.factory()

	t.Equal("12", val)
	t.NoError(err)
}

func (t *BindingTestSuite) TestBindAbstractFactoryError() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	binding.ToFactory(func(a, b Any) (Any, error) {
		return a.(string) + b.(string), nil
	}, 1, 2)
	binding.resolves = NAny{defaultDepVal("1"), defaultDepError(errors.New("error"))}
	val, err := binding.factory()

	t.Nil(val)
	t.Error(err)
}

func (t *BindingTestSuite) TestBindTypedFactory() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	counter := 0
	binding.ToTypedFactory(func(a, b string) (string, error) {
		counter++
		return a + b, nil
	}, 1, 2)
	binding.resolves = NAny{defaultDepVal("1"), defaultDepVal("2")}
	val, err := binding.factory()
	val, err = binding.factory()

	t.Equal("12", val)
	t.NoError(err)
	t.Equal(2, counter)
}

func (t *BindingTestSuite) TestBindTypedFactorySingleton() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	counter := 0
	binding.ToTypedFactory(func(a, b string) (string, error) {
		counter++
		return a + b, nil
	}, 1, 2).InSingletonScope()
	binding.resolves = NAny{defaultDepVal("1"), defaultDepVal("2")}
	val, err := binding.factory()
	val, err = binding.factory()

	t.Equal("12", val)
	t.NoError(err)
	t.Equal(1, counter)
}

func (t *BindingTestSuite) TestBindTypedFactoryNoDeps() {
	binding := &Binding{sync.Once{}, nil, NAny{}, NAny{}}

	binding.ToTypedFactory(func() (string, error) {
		return "12", nil
	})
	val, err := binding.factory()

	t.Equal("12", val)
	t.NoError(err)
}

func TestBindingSuite(t *testing.T) {
	suite.Run(t, new(BindingTestSuite))
}
