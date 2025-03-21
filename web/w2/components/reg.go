package components

import (
	"golang.org/x/sync/errgroup"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
)

func Init() {
	dom.ConsoleLog("Init Components")

	var eg errgroup.Group

	eg.Go(func() error {
		return orbital.Register(OrbitalCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Mod {
			return NewOrbitalComp(di)
		})
	})

	eg.Go(func() error {
		return orbital.Register(TaskbarCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Mod {
			return NewTaskbarComp(di)
		})
	})

	eg.Go(func() error {
		return orbital.Register(TaskbarStartCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Mod {
			return NewTaskbarStartComp(di)
		})
	})

	eg.Go(func() error {
		return orbital.Register(OverlayCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Mod {
			return NewOverlayComp(di)
		})
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
