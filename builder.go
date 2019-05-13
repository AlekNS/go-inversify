package inversify

import "fmt"

func buildContainerImpl(c *containerImpl) {
	for _, bind := range c.factories {
		for inx, dep := range bind.dependencies {
			optdep, hasOpt := dep.(optionalBind)
			if hasOpt {
				dep = optdep.dep
			}

			b, hasDep := c.findFactory(dep)
			if hasDep {
				bind.dependencies[inx] = b.factory
			} else {
				if hasOpt {
					bind.dependencies[inx] = FactoryFunc(func() (Any, error) {
						return nil, nil
					})
				} else {
					panic(fmt.Sprintf("depending %+v not found", dep))
				}
			}
		}
	}
}
