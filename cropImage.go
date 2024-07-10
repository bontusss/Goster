package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

func cropImage(img image.Image) error {
	instWidth, instaheight := 1080, 1080
	linkedWith, linkedHeight := 1200, 628
	twitterWidth, twitterHeight := 1600, 900

	croppedInsta := imaging.Fill(img, instWidth, instaheight, imaging.Center, imaging.Lanczos)
	croppedLinked := imaging.Fill(img, linkedWith, linkedHeight, imaging.Center, imaging.Lanczos)
	croppedTwitter := imaging.Fill(img, twitterWidth, twitterHeight, imaging.Center, imaging.Lanczos)

	linkedOut := "linkedCrop.jpg"
	instaOut := "instaCrop.jpg"
	twitterOut := "twitterCrop.jpg"

	f, err := os.Create(linkedOut)
	if err != nil {
		return err
	}
	defer f.Close()

	k, err := os.Create(instaOut)
	if err != nil {
		return err
	}
	defer k.Close()

	t, err := os.Create(twitterOut)
	if err != nil {
		return err
	}
	defer t.Close()

	err = jpeg.Encode(f, croppedInsta, &jpeg.Options{Quality: 95})
	if err != nil {
		return err
	}

	err = jpeg.Encode(k, croppedLinked, &jpeg.Options{Quality: 95})
	if err != nil {
		return err
	}

	err = jpeg.Encode(t, croppedTwitter, &jpeg.Options{Quality: 95})
	if err != nil {
		return err
	}
	fmt.Println("linked and insta, and twitter images saved")

	return nil
}
