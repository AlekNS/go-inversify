package inversify

import (
	"fmt"
)

func unwrapDependency(dependency Any) (Any, string, bool) {
	// unwrap
	optionalDependency, isOptional := dependency.(optionalBind)
	if isOptional {
		dependency = optionalDependency.dependency
	}

	var name string
	if namedDependency, isNamed := dependency.(namedBind); isNamed {
		dependency = namedDependency.dependency
		name = namedDependency.name
	}

	return dependency, name, isOptional
}

func resolveContainerDependencies(container *containerDefault) error {
	scd := getStronglyConnectedDependencyList(container.factories)
	for _, items := range scd {
		if len(items) > 1 {
			return fmt.Errorf("the container has cycle dependencies: %+v", items)
		}
	}

	for _, bindings := range container.factories {
		for _, binding := range bindings {
			binding.resolves = make(NAny, len(binding.dependencies))

			for dependencyIndex, dependency := range binding.dependencies {
				dependency, name, isOptional := unwrapDependency(dependency)

				if existingBinding, isDependencyFound := container.findFactory(dependency, name); isDependencyFound {
					binding.resolves[dependencyIndex] = existingBinding.factory
				} else {
					if isOptional {
						binding.resolves[dependencyIndex] = FactoryFunc(func() (Any, error) {
							return Any(nil), nil
						})
					} else {
						return fmt.Errorf("dependency %+v is not found", dependency)
					}
				}
			}
		}
	}

	return nil
}
