package dic

import (
	dingoerrors "errors"

	dingodi "github.com/sarulabs/dingo-example/var/lib/services/dic/dependencies/di"
	dingodependencies "github.com/sarulabs/dingo-example/var/lib/services/dic/dependencies"

	garage "github.com/sarulabs/dingo-example/app/models/garage"
	zap "go.uber.org/zap"
	mgov "gopkg.in/mgo.v2"
)

func getDefinitions(provider *dingodependencies.Provider) []dingodi.Definition {
	return []dingodi.Definition{
		{
			Name: "car-manager",
			Scope: "request",
			Build: func(ctn dingodi.Container) (interface{}, error) {


	
				pi1, err := ctn.SafeGet("logger")
				if err != nil {
					return nil, err
				}
	
				p1, ok := pi1.(*zap.Logger)
				if !ok {
					return nil, dingoerrors.New("could not cast parameter Logger of car-manager to *zap.Logger")
				}



	
				pi0, err := ctn.SafeGet("car-repository")
				if err != nil {
					return nil, err
				}
	
				p0, ok := pi0.(*garage.CarRepository)
				if !ok {
					return nil, dingoerrors.New("could not cast parameter Repo of car-manager to *garage.CarRepository")
				}


				return &garage.CarManager{
					Repo: p0,

					Logger: p1,
				}, nil
			},
			Close: func(obj interface{}) {},
		},
		{
			Name: "car-repository",
			Scope: "request",
			Build: func(ctn dingodi.Container) (interface{}, error) {


	
				pi0, err := ctn.SafeGet("mongo")
				if err != nil {
					return nil, err
				}
	
				p0, ok := pi0.(*mgov.Session)
				if !ok {
					return nil, dingoerrors.New("could not cast parameter Session of car-repository to *mgov.Session")
				}


				return &garage.CarRepository{
					Session: p0,
				}, nil
			},
			Close: func(obj interface{}) {},
		},
		{
			Name: "logger",
			Scope: "app",
			Build: func(ctn dingodi.Container) (interface{}, error) {

				d, err := provider.Get("logger")
				if err != nil {
					return nil, err
				}
				b, ok := d.Build.(func() (*zap.Logger, error))
				if !ok {
					return nil, dingoerrors.New("could not cast build of logger to func() (*zap.Logger, error)")
				}

				return b()
			},
			Close: func(obj interface{}) {},
		},
		{
			Name: "mongo",
			Scope: "request",
			Build: func(ctn dingodi.Container) (interface{}, error) {

				d, err := provider.Get("mongo")
				if err != nil {
					return nil, err
				}

	
				pi0, err := ctn.SafeGet("mongo-pool")
				if err != nil {
					return nil, err
				}
	
				p0, ok := pi0.(*mgov.Session)
				if !ok {
					return nil, dingoerrors.New("could not cast parameter 0 of mongo to *mgov.Session")
				}


				b, ok := d.Build.(func(*mgov.Session) (*mgov.Session, error))
				if !ok {
					return nil, dingoerrors.New("could not cast build of mongo to func(*mgov.Session) (*mgov.Session, error)")
				}

				return b(p0)
			},
			Close: func(obj interface{}) {
				d, err := provider.Get("mongo")
				if err != nil {
					return
				}

				c, ok := d.Close.(func(*mgov.Session))
				if !ok {
					return
				}

				o, ok := obj.(*mgov.Session)
				if !ok {
					return
				}

				c(o)
			},
		},
		{
			Name: "mongo-pool",
			Scope: "app",
			Build: func(ctn dingodi.Container) (interface{}, error) {

				d, err := provider.Get("mongo-pool")
				if err != nil {
					return nil, err
				}
				b, ok := d.Build.(func() (*mgov.Session, error))
				if !ok {
					return nil, dingoerrors.New("could not cast build of mongo-pool to func() (*mgov.Session, error)")
				}

				return b()
			},
			Close: func(obj interface{}) {
				d, err := provider.Get("mongo-pool")
				if err != nil {
					return
				}

				c, ok := d.Close.(func(*mgov.Session))
				if !ok {
					return
				}

				o, ok := obj.(*mgov.Session)
				if !ok {
					return
				}

				c(o)
			},
		},
	}
}

