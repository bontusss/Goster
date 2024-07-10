package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
)

// OverlayWithColorCodes overlays the image with squares filled with colors from provided codes
func OverlayWithColorCodes(imgPaths []string, codes []string) ([]image.Image, error) {
	// Check if image paths are nil
	if imgPaths == nil {
		return nil, fmt.Errorf("invalid image paths")
	}

	// Create a slice to hold the loaded images
	imgs := make([]image.Image, len(imgPaths))

	for i, path := range imgPaths {
		img, err := loadImage(path)
		if err != nil {
			return nil, err
		}
		bounds := img.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		// Set rectangle height and calculate width to span the entire image
		rectHeight := 40
		padding := 0
		bottomPadding := 0
		rectWidth := (width - (len(codes)-1)*padding) / len(codes)

		// Create a new image with same dimensions and RGBA format
		newImg := image.NewRGBA(bounds)
		draw.Draw(newImg, bounds, img, image.Point{X: 0, Y: 0}, draw.Src)

		// Loop through color codes and create overlays
		for j, code := range codes {
			color, err := parseColor(code)
			if err != nil {
				return nil, fmt.Errorf("error parsing color code %s: %w", code, err)
			}

			startX := j * (rectWidth + padding)
			endX := startX + rectWidth
			startY := height - rectHeight - bottomPadding
			endY := startY + rectHeight

			// Draw a rectangle filled with the parsed color
			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					newImg.Set(x, y, color)
				}
			}
		}
		imgs[i] = newImg
	}

	return imgs, nil
}

// parseColor parses a color code (hex string) and returns a color.Color
func parseColor(code string) (color.Color, error) {
	if len(code) > 0 && code[0] == '#' {
		code = code[1:] // Remove leading hash (optional)
	}
	var value uint64
	var err error
	value, err = strconv.ParseUint(code, 16, 32) // Parse hex string
	if err != nil {
		return nil, err
	}
	r := uint8((value >> 16) & 0xFF)
	g := uint8((value >> 8) & 0xFF)
	b := uint8(value & 0xFF)
	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// loadImage loads an image from a file path
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening image file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	switch format {
	case "jpeg":
		return img, nil
	case "png":
		return img, nil
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

// saveImage saves an image to a file path in the specified format (JPEG or PNG)
func saveImages(imgs []image.Image, path string, format string) error {
	for i, img := range imgs {
		file, err := os.Create(fmt.Sprintf("%s_%d.%s", path, i, format))
		if err != nil {
			return fmt.Errorf("error creating image file: %w", err)
		}
		defer file.Close()

		switch format {
		case "jpeg":
			err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95}) // Adjust quality as needed
		case "png":
			err = png.Encode(file, img)
		default:
			err = fmt.Errorf("unsupported image format for saving: %s", format)
		}

		if err != nil {
			return fmt.Errorf("error encoding image: %w", err)
		}
	}
	return nil
}

func main() {
	fmt.Println("getting url...")
	photo, err := getRandomImage("https://colosach.onrender.com/api/search", "food", "green")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("url is ", photo.Src.Original)
	fmt.Println("downloading image...")
	err = downloadImage(photo.Src.Original, "./image.jpg")
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("image downloaded")
	
	// Example usage: Replace with your image loading logic
	fmt.Println("loading image")
	img, err := loadImage("image.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	// extract colors deom inage
	fmt.Println("extracting colors")
	colors, err := getColors(img)
	if err != nil {
		fmt.Println("err", err)
	}

	// add avgColor from api to palette
	colors = append(colors, photo.AvgColor)

	// crop to insta and linked sizes
	fmt.Println("starting crop...")
	err = cropImage(img)
	if err != nil {
		fmt.Println("error cropping image")
	}

	// Overlay the image with colors
	fmt.Println("overlaying colors...")
	imagePaths := []string{"instaCrop.jpg", "linkedCrop.jpg", "twitterCrop.jpg"}
	finalImages, err := OverlayWithColorCodes(imagePaths, colors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save the resulting image (replace with your saving logic)
	saveImages(finalImages, "output", "jpeg")

	// adding logo
	fmt.Println("adding logo...")
	i := []string{"output_0.jpeg", "output_1.jpeg", "output_0.jpeg"}
	err = OverlayWithLogoAndSave(i, "./Colos.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	// uploading to instgram
	text := fmt.Sprintf("%s by %s. %s", photo.Alt, photo.Photographer, photo.AvgColor)
	fmt.Println(text)
	// err = uploadToInstagram(text)
	// if err != nil {
	// 	fmt.Println(err)
	// }

}
