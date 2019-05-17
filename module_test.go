package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ModuleTestSuite struct {
	suite.Suite
}

func getModuleTest1() *Module {
	return NewModule("test1").
		Register(func(c ContainerBinder) error {
			c.Bind(1).ToTypedFactory(func(arg string) (string, error) {
				return arg + "hello", nil
			}, 0)
			return nil
		}).
		UnRegister(func(c ContainerBinder) error {
			c.Unbind(1)
			return nil
		})
}

func getModuleTest2() *Module {
	return NewModule("test2").
		Register(func(c ContainerBinder) error {
			c.Bind(2).To(" world!!!")
			return nil
		}).
		UnRegister(func(c ContainerBinder) error {
			c.Unbind(2)
			return nil
		})
}

func (t *ModuleTestSuite) TestBasic() {
	mdl1 := getModuleTest1()
	mdl2 := getModuleTest2()

	c := NewContainer()
	c.Bind(0).ToTypedFactory(func() (string, error) {
		return "Modules, ", nil
	})
	c.Bind("result").ToTypedFactory(func(a, b string) (string, error) {
		return a + b, nil
	}, 1, 2)
	c.Load(mdl1)
	c.Load(mdl2)
	c.Build()

	val, err := c.Get("result")

	t.Equal("Modules, hello world!!!", val)
	t.NoError(err)

	c.UnLoad(mdl1)
	c.UnLoad(mdl2)
	c.Unbind("result")
	c.Build()

	t.False(c.IsBound(1))
}

func TestModuleSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}
