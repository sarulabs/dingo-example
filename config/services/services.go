package services

import (
	"os"
	"time"

	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo-example/app/models/garage"
	"github.com/sarulabs/dingo-example/config/logging"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// Services contains the definitions of the application services.
var Services = []dingo.Def{
	{
		Name:  "logger",
		Scope: dingo.App,
		Build: func() (*zap.Logger, error) {
			return logging.Logger, nil
		},
	},
	{
		Name:  "mongo-pool",
		Scope: dingo.App,
		Build: func() (*mgo.Session, error) {
			return mgo.DialWithTimeout(os.Getenv("MONGO_URL"), 5*time.Second)
		},
		Close: func(s *mgo.Session) {
			s.Close()
		},
		NotForAutoFill: true,
	},
	{
		Name:  "mongo",
		Scope: dingo.Request,
		Build: func(pool *mgo.Session) (*mgo.Session, error) {
			return pool.Copy(), nil
		},
		Params: dingo.Params{
			"0": dingo.Service("mongo-pool"),
		},
		Close: func(s *mgo.Session) {
			s.Close()
		},
	},
	{
		Name:  "car-repository",
		Scope: dingo.Request,
		Build: (*garage.CarRepository)(nil),
	},
	{
		Name:  "car-manager",
		Scope: dingo.Request,
		Build: (*garage.CarManager)(nil),
	},
}
