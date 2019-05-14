package inversify

import (
	"fmt"
)

func resolveContainerDependencies(c *containerDefault) {
	for _, bind := range c.factories {
		// if bind.resolves != nil {
		// 	continue
		// }
		bind.resolves = make(NAny, len(bind.dependencies))
		for inx, dep := range bind.dependencies {
			optdep, hasOpt := dep.(optionalBind)
			if hasOpt {
				dep = optdep.dep
			}

			b, hasDep := c.findFactory(dep)
			if hasDep {
				bind.resolves[inx] = b.factory
			} else {
				if hasOpt {
					bind.resolves[inx] = FactoryFunc(func() (Any, error) {
						return nil, nil
					})
				} else {
					panic(fmt.Sprintf("depending %+v not found", dep))
				}
			}
		}
	}
}
