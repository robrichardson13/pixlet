package gif_api

import (
	"fmt"
	"image"
	"log"

	"github.com/robrichardson13/pixlet/encode"
	"github.com/robrichardson13/pixlet/runtime"
)

func GIF(src []byte, config map[string]string) string {
	script := "temp.star"

	runtime.InitCache(runtime.NewInMemoryCache())

	applet := runtime.Applet{}
	err := applet.Load(script, src, nil)
	if err != nil {
		fmt.Printf("failed to load applet: %v\n", err)
		return "failed to load applet"
	}

	roots, err := applet.Run(config)
	if err != nil {
		log.Printf("Error running script: %s\n", err)
		return "Error running script"
	}
	screens := encode.ScreensFromRoots(roots)

	filter := func(input image.Image) (image.Image, error) {
		return input, nil
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
