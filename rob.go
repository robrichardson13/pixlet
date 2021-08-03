package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"strings"

	"github.com/robrichardson13/pixlet/encode"
	"github.com/robrichardson13/pixlet/runtime"
)

func rob() string {
	script := "message.star"

	if !strings.HasSuffix(script, ".star") {
		fmt.Printf("script file must have suffix .star: %s\n", script)
		return "script file must have suffix .star"
	}

	src, err := ioutil.ReadFile(script)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", script, err)
		return "failed to read file"
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	applet := runtime.Applet{}
	err = applet.Load(script, src, nil)
	if err != nil {
		fmt.Printf("failed to load applet: %v\n", err)
		return "failed to load applet"
	}

	config := map[string]string{}
	roots, err := applet.Run(config)
	if err != nil {
		log.Printf("Error running script: %s\n", err)
		return "Error running script"
	}
	screens := encode.ScreensFromRoots(roots)

	filter := func(input image.Image) (image.Image, error) {
		if magnify <= 1 {
			return input, nil
		}
		in, ok := input.(*image.RGBA)
		if !ok {
			return nil, fmt.Errorf("image not RGBA, very weird")
		}

		out := image.NewRGBA(
			image.Rect(
				0, 0,
				in.Bounds().Dx()*magnify,
				in.Bounds().Dy()*magnify),
		)
		for x := 0; x < in.Bounds().Dx(); x++ {
			for y := 0; y < in.Bounds().Dy(); y++ {
				for xx := 0; xx < 10; xx++ {
					for yy := 0; yy < 10; yy++ {
						out.SetRGBA(
							x*magnify+xx,
							y*magnify+yy,
							in.RGBAAt(x, y),
						)
					}
				}
			}
		}

		return out, nil
	}

	var buf []byte
	buf, err = screens.EncodeGIF(filter)
	if err != nil {
		fmt.Printf("Error rendering: %s\n", err)
		return "Error rendering"
	}

	out := screens.Base64Encode(buf)
	return out
}
