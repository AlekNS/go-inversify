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
	Values map[string]interface{} `inversify:"strkey:values"`
	Value  int                    `inversify:"intkey:1,optional"`

	Config *config `inversify:"inject"`

	TaskRepository iTaskRepository `inversify:"named:gorm"`
	Scheduler      iScheduler      `inversify:"optional"`
}

type AutowireStructTestSuite struct {
	suite.Suite
}

func (t *AutowireStructTestSuite) TestBasic() {
	c := NewContainer()
	c.Bind("values").To(map[string]interface{}{
		"value1": "1",
		"value2": "2",
	})
	c.Bind((*config)(nil)).To(&config{1})
	c.Bind((*iTaskRepository)(nil)).ToFactory(func() (iTaskRepository, error) {
		return &iTaskRepositoryImpl{
			val: 2,
		}, nil
	})
	c.Build()

	s := autowireTestStruct{}
	AutowireStruct(&s)

	// t.NotNil(s.Values)
	// t.NotNil(s.TaskRepository)
	// t.NotNil(s.Scheduler)
	// t.Nil(s.Config)
}

func TestAutowireStructSuite(t *testing.T) {
	suite.Run(t, new(AutowireStructTestSuite))
}
