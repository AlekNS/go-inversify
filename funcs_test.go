package inversify

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FuncsTestSuite struct {
	suite.Suite
}

func wrapTestArgs0() (Any, error)                        { return 0, nil }
func wrapTestArgs1(Any) (Any, error)                     { return 0, nil }
func wrapTestArgs2(Any, Any) (Any, error)                { return 0, nil }
func wrapTestArgs3(Any, Any, Any) (Any, error)           { return 0, nil }
func wrapTestArgs4(Any, Any, Any, Any) (Any, error)      { return 0, nil }
func wrapTestArgs5(Any, Any, Any, Any, Any) (Any, error) { return 0, nil }

func wrapTestArgs6(Any, Any, Any, Any, Any, Any) (Any, error)                      { return 0, nil }
func wrapTestArgs7(Any, Any, Any, Any, Any, Any, Any) (Any, error)                 { return 0, nil }
func wrapTestArgs8(Any, Any, Any, Any, Any, Any, Any, Any) (Any, error)            { return 0, nil }
func wrapTestArgs9(Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error)       { return 0, nil }
func wrapTestArgs10(Any, Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error) { return 0, nil }
func wrapTestArgs10WithError(Any, Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error) {
	return 0, errors.New("")
}

func wrapTestArgs11(Any, Any, Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error) { return 0, nil }
func wrapTestArgs11WithError(Any, Any, Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error) {
	return 0, errors.New("")
}

func (t *FuncsTestSuite) TestWrapPassNotAFunction() {
	t.PanicsWithValue("not a function", func() {
		wrapAbstractApplyFuncAsSlice("invalid value")
	})
}

func (t *FuncsTestSuite) TestWrapFunc0Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs0)
	result, err := method([]Any{})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc1Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs1)
	result, err := method([]Any{1})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc1SliceCustom() {
	method := wrapCustomApplyFuncAsSlice(func (int) (int, error) {
		return 0, nil
	}, func(raw Any) func([]Any) (Any, error) {
		return func(args []Any) (Any, error) {
			return raw.(func(int) (int, error))(args[0].(int))
		}
	})
	result, err := method([]Any{1})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc2Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs2)
	result, err := method([]Any{1, 2})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc3Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs3)
	result, err := method([]Any{1, 2, 3})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc4Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs4)
	result, err := method([]Any{1, 2, 3, 4})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc5Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs5)
	result, err := method([]Any{1, 2, 3, 4, 5})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc6Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs6)
	result, err := method([]Any{1, 2, 3, 4, 5, 6})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc7Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs7)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc8Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs8)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc9Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs9)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8, 9})

	t.Equal(0, result)
	t.NoError(err)
}
func (t *FuncsTestSuite) TestWrapFunc10Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs10)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc10SliceWithError() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs10WithError)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	t.Equal(0, result)
	t.Error(err)
}

func (t *FuncsTestSuite) TestWrapFunc11Slice() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs11)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})

	t.Equal(0, result)
	t.NoError(err)
}

func (t *FuncsTestSuite) TestWrapFunc11SliceWithError() {
	method := wrapAbstractApplyFuncAsSlice(wrapTestArgs11WithError)
	result, err := method([]Any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})

	t.Equal(0, result)
	t.Error(err)
}

func typedFuncForTest(a, b, c, d, e, f string) (string, error) {
	return a + b + c + d + e + f, nil
}

func abstractFuncForTest(a, b, c, d, e, f Any) (Any, error) {
	return a.(string) + b.(string) + c.(string) +
		d.(string) + e.(string) + f.(string), nil
}

func (t *FuncsTestSuite) TestWrapFuncTypedSlice() {
	method := wrapTypedApplyFuncAsSlice(typedFuncForTest)
	result, err := method([]Any{"1", "2", "3", "4", "5", "6"})

	t.Equal("123456", result)
	t.NoError(err)
}

func TestFuncsSuite(t *testing.T) {
	suite.Run(t, new(FuncsTestSuite))
}

func BenchmarkDirectCall(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = typedFuncForTest("1", "2", "3", "4", "5", "6")
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIndirectTypedCall(b *testing.B) {
	method := wrapTypedApplyFuncAsSlice(typedFuncForTest)
	args := []Any{"1", "2", "3", "4", "5", "6"}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = method(args)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIndirectAbstractCall(b *testing.B) {
	method := wrapAbstractApplyFuncAsSlice(abstractFuncForTest)
	args := []Any{"1", "2", "3", "4", "5", "6"}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = method(args)
	}
	b.StopTimer()
	b.ReportAllocs()
}
