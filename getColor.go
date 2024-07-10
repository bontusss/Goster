package main

import (
	"fmt"
	"image"

	"github.com/Edlinorg/prominentcolor"
)

func getColors(img image.Image) ([]string, error) {
	var colorPalette []string
	// file, err := os.Open(fileInput)
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	// img, _, err := image.Decode(file)
	// if err != nil {
	// 	return nil, err
	// }
	colors, err := prominentcolor.Kmeans(img)
	if err != nil {
		return nil, err
	}

	for _, k :=range colors {
		colorPalette = append(colorPalette, fmt.Sprintf("#%02x%02x%02x", k.Color.R, k.Color.G, k.Color.B))
	}
	return colorPalette, nil
}
