package native

import (
	"bufio"
	"bytes"
	"errors"
	_ "fmt"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiifsource "github.com/go-iiif/go-iiif/source"
	"github.com/whosonfirst/go-whosonfirst-mimetypes"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	_ "log"
)

type NativeImage struct {
	iiifimage.Image
	config    *iiifconfig.Config
	source    iiifsource.Source
	source_id string
	id        string
	img       image.Image
	format    string
}

type NativeDimensions struct {
	iiifimage.Dimensions
	bounds image.Rectangle
}

func (d *NativeDimensions) Height() int {
	return d.bounds.Max.Y
}

func (d *NativeDimensions) Width() int {
	return d.bounds.Max.X
}

func (im *NativeImage) Update(body []byte) error {

	buf := bytes.NewBuffer(body)

	img, fmt, err := image.Decode(buf)

	if err != nil {
		return err
	}

	im.img = img
	im.format = fmt

	return nil
}

func (im *NativeImage) Body() []byte {

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	switch im.format {
	case "bmp":
		bmp.Encode(wr, im.img)
	case "jpeg":
		opts := jpeg.Options{Quality: 100}
		jpeg.Encode(wr, im.img, &opts)
	case "png":
		png.Encode(wr, im.img)
	case "gif":
		opts := gif.Options{}
		gif.Encode(wr, im.img, &opts)
	case "tiff":
		opts := tiff.Options{}
		tiff.Encode(wr, im.img, &opts)
	default:
		//
	}

	wr.Flush()
	return b.Bytes()
}

func (im *NativeImage) Format() string {

	return im.format
}

func (im *NativeImage) ContentType() string {

	format := im.Format()

	t := mimetypes.TypesByExtension(format)
	return t[0]
}

func (im *NativeImage) Identifier() string {
	return im.id
}

func (im *NativeImage) Rename(id string) error {
	im.id = id
	return nil
}

func (im *NativeImage) Dimensions() (iiifimage.Dimensions, error) {

	dims := &NativeDimensions{
		bounds: im.img.Bounds(),
	}

	return dims, nil
}

// https://godoc.org/github.com/h2non/img#Options

func (im *NativeImage) Transform(t *iiifimage.Transformation) error {

	return errors.New("Please write me")

	// PLEASE PUT THIS IN A COMMON PACKAGE

	// None of what follows is part of the IIIF spec so it's not clear
	// to me yet how to make this in to a sane interface. For the time
	// being since there is only lipvips we'll just take the opportunity
	// to think about it... (20160917/thisisaaronland)

	// Also note the way we are diligently setting in `im.isgif` in each
	// of the features below. That's because this is a img/libvips-ism
	// and we assume that any of these can encode GIFs because pure-Go and
	// the rest of the code does need to know about it...
	// (20160922/thisisaaronland)

	/*

		if t.Quality == "dither" {

			err = DitherImage(im)

			if err != nil {
				return err
			}

			if fi.Format == "gif" {
				im.isgif = true
			}

		} else if strings.HasPrefix(t.Quality, "primitive:") {

			parts := strings.Split(t.Quality, ":")
			parts = strings.Split(parts[1], ",")

			mode, err := strconv.Atoi(parts[0])

			if err != nil {
				return err
			}

			iters, err := strconv.Atoi(parts[1])

			if err != nil {
				return err
			}

			max_iters := im.config.Primitive.MaxIterations

			if max_iters > 0 && iters > max_iters {
				return errors.New("Invalid primitive iterations")
			}

			alpha, err := strconv.Atoi(parts[2])

			if err != nil {
				return err
			}

			if alpha > 255 {
				return errors.New("Invalid primitive alpha")
			}

			animated := false

			if fi.Format == "gif" {
				animated = true
			}

			opts := PrimitiveOptions{
				Alpha:      alpha,
				Mode:       mode,
				Iterations: iters,
				Size:       0,
				Animated:   animated,
			}

			err = PrimitiveImage(im, opts)

			if err != nil {
				return err
			}

			if fi.Format == "gif" {
				im.isgif = true
			}
		}

	*/

	// END OF none of what follows is part of the IIIF spec

	return nil
}