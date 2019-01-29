package process

import (
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type URI interface {
	URL() string
	String() string
}

type Label string

type Processor interface {
	ProcessURIWithInstructions(URI, string, IIIFInstructions) (URI, iiifimage.Image, error)
}
