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

func main() {
	inputFile, outputDir := parseFlags()
	startTime := time.Now()
	fmt.Println("Программа выполняется...")

	processInputFile(*inputFile, *outputDir)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения программы: %v\n", duration)
}

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

func parseFlags() (*string, *string) {
	inputFile := flag.String("input", "./url.txt", "Путь к файлу с URL-адресами")
	outputDir := flag.String("output", "./output", "Путь к директории для сохранения HTML-файлов")
	flag.Parse()
	setupOutputDir(*outputDir)
	return inputFile, outputDir
}

func setupOutputDir(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка создания директории", err)
		os.Exit(1)
	}
}

func processInputFile(filePath string, outputDir string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ошибка открытия файла:", err)
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		url := scanner.Text()
		url = strings.TrimSpace(url)
		if url != "" {
			if isValidURL(url) {
				fmt.Printf("Valid URL: %s\n", url)
				wg.Add(1)
				go func(i int) { // Используем замыкание для передачи копии i в горутину
					processURL(url, &wg, i, outputDir)
				}(i)
			} else {
				fmt.Printf("Invalid URL: %s\n", url)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		os.Exit(1)
	}

	wg.Wait()
}

func processURL(url string, wg *sync.WaitGroup, i int, outputDir string) {
	defer wg.Done()
	fmt.Printf("Горутина для %s начата\n", url)
	html, err := fetchHTML(url)
	if err != nil {
		fmt.Printf("Ошибка получения HTML для %s: %v\n", url, err)
		return
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("url_%d.html", i+1))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("ошибка создания файла для %s: %v\n", url, err)
		return
	}
	defer file.Close()

	_, err = file.Write(html)
	if err != nil {
		fmt.Printf("Ошибка записи HTML %s: %v\n", url, err)
	}
	fmt.Printf("Горутина для %s завершена\n", url)
}
