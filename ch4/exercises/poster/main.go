package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type MovieInfo struct {
	Title    string
	Response string
	Poster   string
}

func getUrlBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	return resp.Body, nil
}

func downloadPoster(key string, terms []string) error {
	query := url.QueryEscape(strings.Join(terms, " "))
	url := fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&t=%s", key, query)
	body, err := getUrlBody(url)
	if err != nil {
		return err
	}
	defer body.Close()

	var movieInfo MovieInfo
	if err := json.NewDecoder(body).Decode(&movieInfo); err != nil {
		return err
	}

	if movieInfo.Response == "False" {
		return fmt.Errorf("movie not found: %s", strings.Join(terms, " "))
	}

	posterBody, err := getUrlBody(movieInfo.Poster)
	if err != nil {
		return err
	}
	defer posterBody.Close()

	file, err := os.Create(movieInfo.Title)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, posterBody); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("too few arguments")
	}
	err := downloadPoster(os.Args[1], os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}
