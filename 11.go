package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func isValidURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func fetchHTML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body of %s: %v", url, err)
	}

	return body, nil
}

func main() {
	var wg sync.WaitGroup

	inputFile := flag.String("input", "url.txt", "Путь к файлу с URL-адресами")
	outputDir := flag.String("output", "output", "Путь к директории для сохранения HTML-файлов")

	flag.Parse()

	startTime := time.Now()
	fmt.Println("Программа выполняется...")

	err := os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	work := func(url string, i int) {
		defer wg.Done()
		//fmt.Printf("Горутина %d начала выполнение \n", i)
		html, err := fetchHTML(url)
		if err != nil {
			fmt.Printf("Error fetching HTML for %s: %v\n", url, err)
			return
		}

		filename := filepath.Join(*outputDir, fmt.Sprintf("url_%d.html", i+1))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file for %s: %v\n", url, err)
			return
		}
		defer file.Close()

		_, err = file.Write(html)
		if err != nil {
			fmt.Printf("Error writing HTML for %s: %v\n", url, err)
		}

		//fmt.Printf("Горутина %d завершила выполнение \n", i)
	}

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		lineStr := scanner.Text()
		lineStr = strings.TrimSpace(lineStr)
		if lineStr != "" && isValidURL(lineStr) {
			wg.Add(1)
			go work(lineStr, i)
			i++
			fmt.Printf("Valid URL: %s\n", lineStr)
		} else if lineStr != "" {
			fmt.Printf("Invalid URL: %s\n", lineStr)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	wg.Wait()
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Время выполнения программы: %v\n", duration)
}
