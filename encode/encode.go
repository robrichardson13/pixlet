package encode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/gif"

	"github.com/pkg/errors"
	"github.com/robrichardson13/pixlet/render"
)

const (
	WebPKMin                 = 0
	WebPKMax                 = 0
	DefaultScreenDelayMillis = 50
)

type Screens struct {
	roots  []render.Root
	images []image.Image
	delay  int32
}

type ImageFilter func(image.Image) (image.Image, error)

func ScreensFromRoots(roots []render.Root) *Screens {
	screens := Screens{
		roots: roots,
		delay: DefaultScreenDelayMillis,
	}
	if len(roots) > 0 {
		if roots[0].Delay > 0 {
			screens.delay = roots[0].Delay
		}
	}
	return &screens
}

func ScreensFromImages(images ...image.Image) *Screens {
	screens := Screens{
		images: images,
		delay:  DefaultScreenDelayMillis,
	}
	return &screens
}

// Renders a screen to GIF. Optionally pass filters for postprocessing
// each individual frame.
func (s *Screens) EncodeGIF(filters ...ImageFilter) ([]byte, error) {
	images, err := s.render(filters...)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return []byte{}, nil
	}

	g := &gif.GIF{}

	for imIdx, im := range images {
		imRGBA, ok := im.(*image.RGBA)
		if !ok {
			return nil, fmt.Errorf("image %d is %T, require RGBA", imIdx, im)
		}

		palette := color.Palette{}
		idxByColor := map[color.RGBA]int{}

		// Create the palette
		for x := 0; x < imRGBA.Bounds().Dx(); x++ {
			for y := 0; y < imRGBA.Bounds().Dy(); y++ {
				c := imRGBA.RGBAAt(x, y)
				if _, found := idxByColor[c]; !found {
					idxByColor[c] = len(palette)
					palette = append(palette, c)
				}
			}
		}
		if len(palette) > 256 {
			return nil, fmt.Errorf(
				"require <=256 colors, found %d in image %d",
				len(palette), imIdx,
			)
		}

		// Construct the paletted image
		imPaletted := image.NewPaletted(imRGBA.Bounds(), palette)
		for x := 0; x < imRGBA.Bounds().Dx(); x++ {
			for y := 0; y < imRGBA.Bounds().Dy(); y++ {
				imPaletted.SetColorIndex(x, y, uint8(idxByColor[imRGBA.RGBAAt(x, y)]))
			}
		}

		g.Image = append(g.Image, imPaletted)
		g.Delay = append(g.Delay, int(s.delay/10)) // in 100ths of a second
	}

	buf := &bytes.Buffer{}
	err = gif.EncodeAll(buf, g)
	if err != nil {
		return nil, errors.Wrap(err, "encoding")
	}

	return buf.Bytes(), nil
}

func (s *Screens) render(filters ...ImageFilter) ([]image.Image, error) {
	if s.images == nil {
		s.images = render.PaintRoots(true, s.roots...)
	}

	if len(s.images) == 0 {
		return nil, nil
	}

	images := s.images

	if len(filters) > 0 {
		images = []image.Image{}
		for _, im := range s.images {
			for _, f := range filters {
				imFiltered, err := f(im)
				if err != nil {
					return nil, err
				}
				im = imFiltered
			}
			images = append(images, im)
		}
	}

	return images, nil
}

func (s *Screens) Base64Encode(message []byte) string {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return string(b)
}
