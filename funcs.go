package inversify

import (
	"reflect"
)

var dependencyNotFound interface{}

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

	reflectedValue := reflect.ValueOf(f)
	reflectedType := reflectedValue.Type()
	if reflectedType.Kind() != reflect.Func {
		panic("not a function")
	}

	argCount := reflectedType.NumIn()

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

	argValue := make([]reflect.Value, argCount, argCount)

	return func(args []Any) (Any, error) {
		for index, argument := range args {
			if argument == nil {
				argValue[index] = reflect.ValueOf(&dependencyNotFound).Elem()
			} else {
				argValue[index] = reflect.ValueOf(argument)
			}
		}

		results := reflectedValue.Call(argValue)

		if results[1].IsNil() {
			return results[0].Interface(), nil
		}
		if results[0].IsNil() {
			return nil, results[1].Interface().(error)
		}
		return results[0].Interface(), results[1].Interface().(error)
	}
}
