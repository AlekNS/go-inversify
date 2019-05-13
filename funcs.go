package inversify

import (
	"reflect"
)

func wrapApplyFunc0AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func() (Any, error))()
	}
}

func wrapApplyFunc1AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any) (Any, error))(
			args[0])
	}
}

func wrapApplyFunc2AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any) (Any, error))(
			args[0], args[1])
	}
}

func wrapApplyFunc3AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any) (Any, error))(
			args[0], args[1], args[2])
	}
}

func wrapApplyFunc4AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3])
	}
}

func wrapApplyFunc5AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4])
	}
}

func wrapApplyFunc6AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4], args[5])
	}
}

func wrapApplyFunc7AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4], args[5], args[6])
	}
}

func wrapApplyFunc8AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7])
	}
}

func wrapApplyFunc9AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	}
}

func wrapApplyFunc10AsSlice(funcRaw anyFunc) func([]Any) (Any, error) {
	return func(args []Any) (Any, error) {
		return funcRaw.(func(Any, Any, Any, Any, Any, Any, Any, Any, Any, Any) (Any, error))(
			args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9])
	}
}

func wrapTypedApplyFuncAsSlice(f anyFunc) func([]Any) (Any, error) {
	return wrapApplyFuncAsSlice(f, nil, false)
}

func wrapAbstractApplyFuncAsSlice(f anyFunc) func([]Any) (Any, error) {
	return wrapApplyFuncAsSlice(f, nil, true)
}

func wrapCustomApplyFuncAsSlice(f anyFunc, customConv customTranslator) func([]Any) (Any, error) {
	return wrapApplyFuncAsSlice(f, nil, false)
}

func wrapApplyFuncAsSlice(f anyFunc, customConv Any, isAbstract bool) func([]Any) (Any, error) {
	if customConv != nil {
		return customConv.(customTranslator)(f)
	}

	reflVal := reflect.ValueOf(f)
	reflType := reflVal.Type()
	if reflType.Kind() != reflect.Func {
		panic("not a function")
	}

	argCount := reflType.NumIn()
	argValue := make([]reflect.Value, argCount, argCount)

	if isAbstract {
		switch argCount {
		case 0:
			return wrapApplyFunc0AsSlice(f)
		case 1:
			return wrapApplyFunc1AsSlice(f)
		case 2:
			return wrapApplyFunc2AsSlice(f)
		case 3:
			return wrapApplyFunc3AsSlice(f)
		case 4:
			return wrapApplyFunc4AsSlice(f)
		case 5:
			return wrapApplyFunc5AsSlice(f)
		case 6:
			return wrapApplyFunc6AsSlice(f)
		case 7:
			return wrapApplyFunc7AsSlice(f)
		case 8:
			return wrapApplyFunc8AsSlice(f)
		case 9:
			return wrapApplyFunc9AsSlice(f)
		case 10:
			return wrapApplyFunc10AsSlice(f)
		default:
		}
	}

	return func(args []Any) (Any, error) {
		for inx, argument := range args {
			argValue[inx] = reflect.ValueOf(argument)
		}

		results := reflVal.Call(argValue)

		if results[1].IsNil() {
			return results[0].Interface(), nil
		}

		return results[0].Interface(), results[1].Interface().(error)
	}
}
