package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

// OverlayWithLogoAndSave overlays the image with a PNG logo on the top right corner and saves the images
func OverlayWithLogoAndSave(imgPaths []string, logoPath string) error {
	// Check if image is nil
	if imgPaths == nil {
		return fmt.Errorf("invalid image")
	}

	// Load the logo image
	logo, err := loadImage(logoPath)
	if err != nil {
		return err
	}

	for i, path := range imgPaths {
		img, err := loadImage(path)
		if err != nil {
			return err
		}
		bounds := img.Bounds()

		// Create a new image with same dimensions and RGBA format
		newImg := image.NewRGBA(bounds)
		draw.Draw(newImg, bounds, img, image.Point{X: 0, Y: 0}, draw.Src)

		// Calculate the position to place the logo
		logoBounds := logo.Bounds()
		logoX := bounds.Max.X - logoBounds.Max.X - 20 // 20px padding from the right
		logoY := 20 // 20px padding from the top

		// Draw the logo on the image
		draw.Draw(newImg, logoBounds.Add(image.Point{X: logoX, Y: logoY}), logo, image.Point{X: 0, Y: 0}, draw.Over)

		// Save the resulting image
		file, err := os.Create(fmt.Sprintf("final%d.jpg", i+1))
		if err != nil {
			return fmt.Errorf("error creating image file: %w", err)
		}
		defer file.Close()

		err = jpeg.Encode(file, newImg, &jpeg.Options{Quality: 95}) // Adjust quality as needed
		if err != nil {
			return fmt.Errorf("error encoding image: %w", err)
		}
	}

	return nil
}
