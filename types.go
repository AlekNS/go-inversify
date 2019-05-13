package inversify

type (
	// Any specifies abstract type
	Any = interface{}
	// NAny specifies slice of abstract types
	NAny = []interface{}

	// FactoryFunc defines abstract factory
	FactoryFunc func() (Any, error)

	// anyFunc defines pointer to an abstract function
	anyFunc = interface{}

	// AnyReturnFunc defines pointer to an abstract function that returns only one single value
	AnyReturnFunc = interface{}

	customTranslator = func(anyFunc) func([]Any) (anyFunc, error)
)
