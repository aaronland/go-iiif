package driver

import (
	"errors"
	iiifcache "github.com/go-iiif/go-iiif/cache"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver interface {
	NewImageFromConfigWithSource(*iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(*iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(*iiifconfig.Config, string) (iiifimage.Image, error)
}

func RegisterDriver(name string, driver Driver) {

	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("iiif: Register driver is nil")

	}

	if _, dup := drivers[name]; dup {
		panic("index: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers = make(map[string]Driver)
}

func Drivers() []string {

	driversMu.RLock()
	defer driversMu.RUnlock()

	var list []string

	for name := range drivers {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}

func NewDriver(name string) (Driver, error) {

	driversMu.RLock()
	defer driversMu.RUnlock()

	dr, ok := drivers[name]

	if !ok {
		return nil, errors.New("Invalid driver")
	}

	return dr, nil
}

func NewDriverFromConfig(config *iiifconfig.Config) (Driver, error) {
	return NewDriver(config.Graphics.Source.Name)
}