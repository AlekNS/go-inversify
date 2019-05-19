package inversify

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseAutowireField(fieldTagValue string) (isOptional, isKeyInt bool, named, keyName string) {
	parts := strings.Split(fieldTagValue, ",")
	for _, part := range parts {
		keyvalue := strings.Split(part, ":")
		switch keyvalue[0] {
		case "opt":
			fallthrough
		case "optional":
			isOptional = true
		case "strkey":
			keyName = keyvalue[1]
		case "intkey":
			isKeyInt = true
			keyName = keyvalue[1]
		case "named":
			named = keyvalue[1]
		}
	}
	return
}

// AutowireStruct injects dependencies based on annotations and types
func AutowireStruct(container Container, structure Any) error {
	refVal := reflect.ValueOf(structure)
	if refVal.Kind() != reflect.Struct &&
		refVal.Kind() != reflect.Ptr {
		panic("passed invalid type for the autowire")
	}
	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}
	refType := refVal.Type()

	for fieldIndex := 0; fieldIndex < refVal.NumField(); fieldIndex++ {
		fieldType := refType.Field(fieldIndex)
		fieldConfigStr, exists := fieldType.Tag.Lookup("inversify")
		if !exists {
			continue
		}

		fieldValue := refVal.Field(fieldIndex)
		if !fieldValue.CanSet() {
			panic(fmt.Sprintf(`field "%s" of "%v" couldn't be set`, fieldType.Name, refType))
		}

		isOptional, isKeyInt, named, strKeyName := parseAutowireField(fieldConfigStr)

		dependency := reflectInterfacePointers(fieldValue.Interface())

		if len(strKeyName) > 0 {
			if isKeyInt {
				intKey, err := strconv.Atoi(strKeyName)
				if err != nil {
					panic(fmt.Sprintf(`field "%s" of "%v" has invalid integer key`, fieldType.Name, refType))
				}
				dependency = intKey
			} else {
				dependency = strKeyName
			}
		} else {
			if fieldType.Type.Kind() == reflect.Interface {
				dependency = reflect.New(reflect.PtrTo(fieldValue.Type())).Elem().Interface()
			}
		}

		if !container.IsBound(dependency, named) {
			if isOptional {
				continue
			}
			panic(fmt.Sprintf(`dependency not found for field "%s" of "%v"`, fieldType.Name, refType))
		}

		result, err := container.Get(dependency, named)
		if err != nil {
			return err
		}

		fieldValue.Set(reflect.ValueOf(result))
	}

	return nil
}
