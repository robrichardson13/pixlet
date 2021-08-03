package gif_api

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"strings"

	"github.com/robrichardson13/pixlet/encode"
	"github.com/robrichardson13/pixlet/runtime"
)

func GIF() string {
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
