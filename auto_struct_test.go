package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type config struct{ val int }
type iTaskRepository interface{}
type iTaskRepositoryImpl struct{ val int }
type iScheduler interface{}
type iSchedulerImpl struct{}

// autowireTestStruct .
type autowireTestStruct struct {
	Values1 map[string]interface{} `inversify:"strkey:values1"`
	Values2 string                 `inversify:"strkey:values2,optional"`
	Value1  int                    `inversify:"intkey:1,named:another"`
	Value2  int                    `inversify:"intkey:1,optional"`
	Config  *config                `inversify:""`

	TaskRepository iTaskRepository `inversify:""`
	Scheduler      iScheduler      `inversify:"optional"`
}

type AutowireStructTestSuite struct {
	suite.Suite
}

func (t *AutowireStructTestSuite) TestBasic() {
	c := NewContainer("base")
	c.Bind("values1").To(map[string]interface{}{
		"value1": "1",
		"value2": "2",
	})
	c.Bind((*config)(nil)).To(&config{1})
	c.Bind((*iTaskRepository)(nil)).ToFactory(func() (Any, error) {
		return &iTaskRepositoryImpl{
			val: 2,
		}, nil
	})
	c.Bind(1, "another").To(1000)
	c.Build()

	s := autowireTestStruct{}
	err := AutowireStruct(c, &s)

	t.NoError(err)
	t.NotNil(s.Values1)
	t.Equal("", s.Values2)
	t.Equal(1000, s.Value1)
	t.NotNil(s.TaskRepository)
	t.Nil(s.Scheduler)
	t.NotNil(s.Config)
}

func TestAutowireStructSuite(t *testing.T) {
	suite.Run(t, new(AutowireStructTestSuite))
}

func BenchmarkContainerAutowireStructure(b *testing.B) {
	c := NewContainer("base")
	c.Bind("values1").To(map[string]interface{}{
		"value1": "1",
		"value2": "2",
	})
	c.Bind((*config)(nil)).To(&config{1})
	c.Bind((*iTaskRepository)(nil)).ToFactory(func() (Any, error) {
		return &iTaskRepositoryImpl{
			val: 2,
		}, nil
	})
	c.Bind(1, "another").To(1000)
	c.Build()

	s := autowireTestStruct{}

	b.ReportAllocs()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Config = nil
		s.Scheduler = nil
		s.TaskRepository = nil
		AutowireStruct(c, &s)
	}
	b.StopTimer()
}
