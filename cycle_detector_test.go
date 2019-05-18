package inversify

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CycleDetectorTestSuite struct {
	suite.Suite
	factories map[Any]map[string]*Binding
}

func (t *CycleDetectorTestSuite) BeforeTest(suiteName, testName string) {
	//  a       b      c
	//   ^\   /^^^\   /^
	//      d<--|---e
	//       ^\ |  /^
	//          f
	t.factories = map[Any]map[string]*Binding{
		"a": map[string]*Binding{"": &Binding{dependencies: NAny{}}},
		"b": map[string]*Binding{"": &Binding{dependencies: NAny{}}},
		"c": map[string]*Binding{"": &Binding{dependencies: NAny{}}},
		"d": map[string]*Binding{"": &Binding{dependencies: NAny{"a", "b"}}},
		"e": map[string]*Binding{"": &Binding{dependencies: NAny{"b", "c", "d"}}},
		"f": map[string]*Binding{"": &Binding{dependencies: NAny{"b", "d", "e"}}},
	}
}

func (t *CycleDetectorTestSuite) TestEmpty() {
	factories := map[Any]map[string]*Binding{}
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 0)
}

func (t *CycleDetectorTestSuite) TestNoCycles() {
	factories := t.factories
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 6)
	for i := 0; i < len(scDeps); i++ {
		t.Lenf(scDeps[i], 1, "scDeps[%v]=%v should contains 1", i, len(scDeps[i]))
	}
}

func (t *CycleDetectorTestSuite) TestSelfCycles() {
	factories := t.factories
	factories["b"] = map[string]*Binding{"": &Binding{dependencies: NAny{"b"}}}
	factories["e"] = map[string]*Binding{"": &Binding{dependencies: NAny{"b", "c", "e", "d"}}}
	scDeps := getStronglyConnectedDependencyList(factories)

	t.Len(scDeps, 8)
	for i := 0; i < len(scDeps); i++ {
		t.Contains([]int{1, 2}, len(scDeps[i]))
		if len(scDeps[i]) == 2 {
			t.Contains([]string{"b", "e"}, scDeps[i][0].(PairDependencyName).Symbol)
		}
	}
}

func (t *CycleDetectorTestSuite) TestWithCycles() {
	factories := t.factories
	factories["a"] = map[string]*Binding{"": &Binding{dependencies: NAny{"f"}}}

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
