package dic

import (
	dingoerrors "errors"
	dingolog "log"
	dingohttp "net/http"

	dingo "github.com/sarulabs/dingo"
	dingodi "github.com/sarulabs/dingo-example/var/lib/services/dic/dependencies/di"
	dingodependencies "github.com/sarulabs/dingo-example/var/lib/services/dic/dependencies"

	garage "github.com/sarulabs/dingo-example/app/models/garage"
	zap "go.uber.org/zap"
	mgov "gopkg.in/mgo.v2"
)

// C retrieves a Container from an interface.
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its Context
//   for the dingo.ContainerKey("dingo") key.
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}

	r, ok := i.(*dingohttp.Request)
	if !ok && ErrorCallback != nil {
		ErrorCallback(dingoerrors.New("could not get container with C()"))
		return nil
	}
	
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok && ErrorCallback != nil {
		ErrorCallback(dingoerrors.New("could not get container from *http.Request"))
		return nil
	}

	return c
}

// ErrorCallback is a function that is called
// when there is an error while retrieving an object
// with the Get method (and its derivatives).
// The function can be changed to match the needs of your application.
var ErrorCallback = func(err error) {
	dingolog.Println(err.Error())
}

// NewContainer creates a new Container.
// If no scope is provided, dingo.App, dingo.Request and dingo.SubRequest are used.
// The returned Container has the wider scope (dingo.App).
// The SubContainer() method should be called to get a Container in a narrower scope.
func NewContainer(scopes ...string) (*Container, error) {
	if dingo.Version != "1" {
		return nil, dingoerrors.New("The generated code requires github.com/sarulabs/dingo at version 1, but got version " + dingo.Version)
	}

	if len(scopes) == 0 {
		scopes = []string{dingo.App, dingo.Request, dingo.SubRequest}
	}

	b, err := dingodi.NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}

	provider := dingodependencies.NewProvider()
	if err := provider.Load(); err != nil {
		return nil, err
	}

	for _, d := range getDefinitions(provider) {
		if err := b.AddDefinition(d); err != nil {
			return nil, err
		}
	}

	return &Container{ctn: b.Build()}, nil
}

// Container represents a dependency injection container.
// A Container has a scope and may have a parent with a wider scope
// and children with a narrower scope.
// Objects can be retrieved from the Container.
// If the desired object does not already exist in the Container,
// it is built thanks to the object Definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn dingodi.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes narrower than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next subscope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object needs to belong to this scope or a wider one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) Get(name string) interface{} {
	o, err := c.ctn.SafeGet(name)
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a narrower scope.
// To do so UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to Delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
func (c *Container) UnscopedGet(name string) interface{} {
	o, err := c.ctn.UnscopedSafeGet(name)
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// Clean deletes the sub-container created by UnscopedSafeGet or UnscopedGet.
func (c *Container) Clean() {
	c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls their Close function if it exists.
// It will also call DeleteWithSubContainers on each child Container
// and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
func (c *Container) DeleteWithSubContainers() {
	c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers but do not delete the subcontainers.
// If the Container has subcontainers, it will not be deleted right away.
// The deletion only occurs when all the subcontainers have been deleted.
func (c *Container) Delete() {
	c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetCarManager works like SafeGet but only for CarManager.
// It does not return an interface but a *garage.CarManager.
func (c *Container) SafeGetCarManager() (*garage.CarManager, error) {
	i, err := c.ctn.SafeGet("car-manager")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*garage.CarManager)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *garage.CarManager")
	}

	return o, nil
}

// GetCarManager is similar to SafeGetCarManager but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) GetCarManager() *garage.CarManager {
	o, err := c.SafeGetCarManager()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGetCarManager works like UnscopedSafeGet but only for CarManager.
// It does not return an interface but a *garage.CarManager.
func (c *Container) UnscopedSafeGetCarManager() (*garage.CarManager, error) {
	i, err := c.ctn.UnscopedSafeGet("car-manager")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*garage.CarManager)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *garage.CarManager")
	}

	return o, nil
}

// UnscopedGetCarManager is similar to UnscopedSafeGetCarManager but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGetCarManager() *garage.CarManager {
	o, err := c.UnscopedSafeGetCarManager()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// CarManager is similar to GetCarManager.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCarManager method.
// If the container can not be retrieved, it returns the default value for the returned type.
func CarManager(i interface{}) *garage.CarManager {
	c := C(i)
	if c == nil {
		return nil
	}
	return c.GetCarManager()
}

// SafeGetCarRepository works like SafeGet but only for CarRepository.
// It does not return an interface but a *garage.CarRepository.
func (c *Container) SafeGetCarRepository() (*garage.CarRepository, error) {
	i, err := c.ctn.SafeGet("car-repository")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*garage.CarRepository)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *garage.CarRepository")
	}

	return o, nil
}

// GetCarRepository is similar to SafeGetCarRepository but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) GetCarRepository() *garage.CarRepository {
	o, err := c.SafeGetCarRepository()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGetCarRepository works like UnscopedSafeGet but only for CarRepository.
// It does not return an interface but a *garage.CarRepository.
func (c *Container) UnscopedSafeGetCarRepository() (*garage.CarRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("car-repository")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*garage.CarRepository)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *garage.CarRepository")
	}

	return o, nil
}

// UnscopedGetCarRepository is similar to UnscopedSafeGetCarRepository but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGetCarRepository() *garage.CarRepository {
	o, err := c.UnscopedSafeGetCarRepository()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// CarRepository is similar to GetCarRepository.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetCarRepository method.
// If the container can not be retrieved, it returns the default value for the returned type.
func CarRepository(i interface{}) *garage.CarRepository {
	c := C(i)
	if c == nil {
		return nil
	}
	return c.GetCarRepository()
}

// SafeGetLogger works like SafeGet but only for Logger.
// It does not return an interface but a *zap.Logger.
func (c *Container) SafeGetLogger() (*zap.Logger, error) {
	i, err := c.ctn.SafeGet("logger")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*zap.Logger)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *zap.Logger")
	}

	return o, nil
}

// GetLogger is similar to SafeGetLogger but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) GetLogger() *zap.Logger {
	o, err := c.SafeGetLogger()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGetLogger works like UnscopedSafeGet but only for Logger.
// It does not return an interface but a *zap.Logger.
func (c *Container) UnscopedSafeGetLogger() (*zap.Logger, error) {
	i, err := c.ctn.UnscopedSafeGet("logger")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*zap.Logger)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *zap.Logger")
	}

	return o, nil
}

// UnscopedGetLogger is similar to UnscopedSafeGetLogger but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGetLogger() *zap.Logger {
	o, err := c.UnscopedSafeGetLogger()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// Logger is similar to GetLogger.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetLogger method.
// If the container can not be retrieved, it returns the default value for the returned type.
func Logger(i interface{}) *zap.Logger {
	c := C(i)
	if c == nil {
		return nil
	}
	return c.GetLogger()
}

// SafeGetMongo works like SafeGet but only for Mongo.
// It does not return an interface but a *mgov.Session.
func (c *Container) SafeGetMongo() (*mgov.Session, error) {
	i, err := c.ctn.SafeGet("mongo")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*mgov.Session)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *mgov.Session")
	}

	return o, nil
}

// GetMongo is similar to SafeGetMongo but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) GetMongo() *mgov.Session {
	o, err := c.SafeGetMongo()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGetMongo works like UnscopedSafeGet but only for Mongo.
// It does not return an interface but a *mgov.Session.
func (c *Container) UnscopedSafeGetMongo() (*mgov.Session, error) {
	i, err := c.ctn.UnscopedSafeGet("mongo")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*mgov.Session)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *mgov.Session")
	}

	return o, nil
}

// UnscopedGetMongo is similar to UnscopedSafeGetMongo but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGetMongo() *mgov.Session {
	o, err := c.UnscopedSafeGetMongo()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// Mongo is similar to GetMongo.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetMongo method.
// If the container can not be retrieved, it returns the default value for the returned type.
func Mongo(i interface{}) *mgov.Session {
	c := C(i)
	if c == nil {
		return nil
	}
	return c.GetMongo()
}

// SafeGetMongoPool works like SafeGet but only for MongoPool.
// It does not return an interface but a *mgov.Session.
func (c *Container) SafeGetMongoPool() (*mgov.Session, error) {
	i, err := c.ctn.SafeGet("mongo-pool")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*mgov.Session)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *mgov.Session")
	}

	return o, nil
}

// GetMongoPool is similar to SafeGetMongoPool but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) GetMongoPool() *mgov.Session {
	o, err := c.SafeGetMongoPool()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGetMongoPool works like UnscopedSafeGet but only for MongoPool.
// It does not return an interface but a *mgov.Session.
func (c *Container) UnscopedSafeGetMongoPool() (*mgov.Session, error) {
	i, err := c.ctn.UnscopedSafeGet("mongo-pool")
	if err != nil {
		return nil, err
	}

	o, ok := i.(*mgov.Session)
	if !ok {
		return nil, dingoerrors.New("could not cast object to *mgov.Session")
	}

	return o, nil
}

// UnscopedGetMongoPool is similar to UnscopedSafeGetMongoPool but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGetMongoPool() *mgov.Session {
	o, err := c.UnscopedSafeGetMongoPool()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// MongoPool is similar to GetMongoPool.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetMongoPool method.
// If the container can not be retrieved, it returns the default value for the returned type.
func MongoPool(i interface{}) *mgov.Session {
	c := C(i)
	if c == nil {
		return nil
	}
	return c.GetMongoPool()
}



