package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	query := `[out:json];node(around:15000,-7.9333,112.3083)[amenity];out;`
	apiURL := "http://overpass-api.de/api/interpreter"
	
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader("data="+url.QueryEscape(query)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "LOKARI-Dev/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body)[:min(500, len(body))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
