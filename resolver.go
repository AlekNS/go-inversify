package inversify

import (
	"fmt"
)

func resolveContainerDependencies(container *containerDefault) error {
	for _, binding := range container.factories {
		binding.resolves = make(NAny, len(binding.dependencies))

		for dependencyIndex, dependency := range binding.dependencies {

			optionalDependency, isOptional := dependency.(optionalBind)
			if isOptional {
				dependency = optionalDependency.dependency
			}

			if optionalBinding, isDependencyFound := container.findFactory(dependency); isDependencyFound {
				binding.resolves[dependencyIndex] = optionalBinding.factory
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

	scd := getStronglyConnectedDependencyList(container.factories)
	for _, items := range scd {
		if len(items) > 1 {
			return fmt.Errorf("the container has cycle dependencies: %+v", items)
		}
	}

	return nil
}
