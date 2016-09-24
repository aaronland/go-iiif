package image

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"strings"
)

/*

things I am unsure of (20160911/thisisaaronland)
1. whether this should call sanitize.SanitizeString
2. whether this should even exist in its own package

*/

func ScrubIdentifier(identifier string) (string, error) {

	clean, err := url.QueryUnescape(identifier)

	if err != nil {
		return "", err
	}

	clean = strings.Replace(clean, "../", "", -1)
	return clean, nil
}

func IIIFImageToGolangImage(im Image) (image.Image, error) {

	var goimg image.Image
	var err error

	content_type := im.ContentType()

	if content_type == "image/gif" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = gif.Decode(buf)

	} else if content_type == "image/jpeg" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = jpeg.Decode(buf)

	} else if content_type == "image/png" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = png.Decode(buf)

	} else if content_type == "image/tiff" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = tiff.Decode(buf)

	} else if content_type == "image/webp" {

		buf := bytes.NewBuffer(im.Body())
		goimg, err = webp.Decode(buf)

	} else {
		msg := fmt.Sprintf("Unsupported content type '%s' for decoding", content_type)
		err = errors.New(msg)
	}

	if err != nil {
		return nil, err
	}

	return goimg, nil
}

func GolangImageToIIIFImage(goimg image.Image, im Image) error {

	body, err := GolangImageToBytes(goimg, im.ContentType())

	if err != nil {
		return err
	}

	return im.Update(body)
}

func GolangImageToBytes(goimg image.Image, content_type string) ([]byte, error) {

	var out *bytes.Buffer
	var err error

	if content_type == "image/gif" {

		/*
			opts := gif.Options{
				NumColors: 256,
			}
		*/

		out = new(bytes.Buffer)
		err = gif.Encode(out, goimg, nil)

	} else if content_type == "image/jpeg" {

		out = new(bytes.Buffer)
		err = jpeg.Encode(out, goimg, nil)

	} else if content_type == "image/png" {

		out = new(bytes.Buffer)
		err = png.Encode(out, goimg)

	} else if content_type == "image/tiff" {

		out = new(bytes.Buffer)
		err = tiff.Encode(out, goimg, nil)

	} else {

		msg := fmt.Sprintf("Unsupported content type '%s' for encoding", content_type)
		err = errors.New(msg)
	}

	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}