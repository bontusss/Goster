package main

import (
	"fmt"
	"os"

	"github.com/ahmdrz/goinsta/v2"
)

func uploadToInstagram(caption string) error {
	insta := goinsta.New("colosach_", "09031945894IG@")
	if err := insta.Login(); err != nil {
		return err
	}

	file, err := os.Open("final1.jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	// upload media
	_, err = insta.UploadPhoto(file, caption, 80, 0)
	if err != nil {
		return err
	}

	fmt.Println("image uploaded successfully")
	return nil

}
