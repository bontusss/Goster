package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getRandomImage(apiUrl, query, color string) (*Photo, error) {
	reqBody := map[string]string{
		"query": query,
		"color": color,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned status code %d", resp.StatusCode)
	}

	var apiResponse PhotoResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Println("error1", err)
		return nil, err
	}

	if len(apiResponse.Photos) == 0 {
		return nil, errors.New("np images found in response")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randonImage :=apiResponse.Photos[rand.Intn(len(apiResponse.Photos))]
	return &randonImage, nil
}

func downloadImage(url, filePath string) error {
	fmt.Println("fetching image...")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("saving image...")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}