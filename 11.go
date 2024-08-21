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
	"time"
)

func isValidURL(url string) bool {
	resp, err := http.Head(url)// Отправляет HTTP HEAD запрос к указанному URL.
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func fetchHTML(url string) ([]byte, error) {// Извлекает HTML содержимое из указанного URL.
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
	// Определение флагов
	inputFile := flag.String("input", "url.txt", "Путь к файлу с URL-адресами")
	outputDir := flag.String("output", "output", "Путь к директории для сохранения HTML-файлов")

	// Анализ аргументов командной строки
	flag.Parse()

	startTime := time.Now() // Получаем время начала выполнения программы
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

	scanner := bufio.NewScanner(file)\\чтение файла построчно
	var validUrls []string
	for scanner.Scan() {
		lineStr := scanner.Text()
		lineStr = strings.TrimSpace(lineStr)
		if lineStr != "" && isValidURL(lineStr) {
			validUrls = append(validUrls, lineStr)
			fmt.Printf("Valid URL: %s\n", lineStr)
		} else if lineStr != "" {
			fmt.Printf("Invalid URL: %s\n", lineStr)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	for i, url := range validUrls {
		html, err := fetchHTML(url)
		if err != nil {
			fmt.Printf("Error fetching HTML for %s: %v\n", url, err)
			continue
		}

		filename := filepath.Join(*outputDir, fmt.Sprintf("url_%d.html", i+1))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file for %s: %v\n", url, err)
			continue
		}
		defer file.Close()

		_, err = file.Write(html)
		if err != nil {
			fmt.Printf("Error writing HTML for %s: %v\n", url, err)
		}
	}

	endTime := time.Now()              // Получаем время окончания выполнения программы
	duration := endTime.Sub(startTime) // Вычисляем длительность выполнения программы
	fmt.Printf("Время выполнения программы: %v\n", duration)
}
