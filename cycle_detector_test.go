package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CycleDetectorTestSuite struct {
	suite.Suite
}

func (t *CycleDetectorTestSuite) getFactories() map[Any]*Binding {
	return map[Any]*Binding{
		"a": &Binding{dependencies: NAny{}},
		"b": &Binding{dependencies: NAny{}},
		"c": &Binding{dependencies: NAny{}},
		"d": &Binding{dependencies: NAny{"a", "b"}},
		"e": &Binding{dependencies: NAny{"b", "c", "d"}},
		"f": &Binding{dependencies: NAny{"b", "d", "e"}},
	}
}

func (t *CycleDetectorTestSuite) TestEmpty() {
	factories := map[Any]*Binding{}
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 0)
}

func (t *CycleDetectorTestSuite) TestNoCycles() {
	factories := t.getFactories()
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 6)
	for i := 0; i < len(scDeps); i++ {
		t.Lenf(scDeps[i], 1, "scDeps[%v]=%v should contains 1", i, len(scDeps[i]))
	}
}

func (t *CycleDetectorTestSuite) TestSelfCycles() {
	factories := t.getFactories()
	factories["b"] = &Binding{dependencies: NAny{"b"}}
	factories["e"] = &Binding{dependencies: NAny{"b", "c", "e", "d"}}
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 8)
	for i := 0; i < len(scDeps); i++ {
		t.Contains([]int{1, 2}, len(scDeps[i]))
		if len(scDeps[i]) == 2 {
			t.Contains([]string{"b", "e"}, scDeps[i][0].(string))
		}
	}
}

func (t *CycleDetectorTestSuite) TestWithCycles() {
	factories := t.getFactories()
	factories["a"] = &Binding{dependencies: NAny{"f"}}

	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 3)
	isStrongDepsDetected := false
	for i := 0; i < len(scDeps); i++ {
		if len(scDeps[i]) == 4 {
			isStrongDepsDetected = true
		}
	}
	t.True(isStrongDepsDetected)
}

func TestCycleDetectorSuite(t *testing.T) {
	suite.Run(t, new(CycleDetectorTestSuite))
}
