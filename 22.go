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
	inputFile, outputDir := parseFlags()// Парсим флаги командной строки для определения входного файла и директории вывода
	startTime := time.Now() // время начала программы
	fmt.Println("Программа выполняется...")

	processInputFile(*inputFile, *outputDir)// Обрабатываем входной файл

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения программы: %v\n", duration)
}

func fetchHTML(url string) ([]byte, error) {// Функция для получения HTML-кода страницы по URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения %s: %v", url, err)// Возвращение ошибки, если запрос к URL не удался
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения %s: %v", url, err)// Возвращение ошибки, если чтение тела ответа не удалось
	}

	return body, nil
}

func parseFlags() (*string, *string) {// Функция для разбора флагов командной строки
	inputFile := flag.String("input", "./url.txt", "Путь к файлу с URL-адресами")
	outputDir := flag.String("output", "./output", "Путь к директории для сохранения HTML-файлов")
	//flag.PrintDefaults()
	flag.Parse()
	setupOutputDir(*outputDir)
	
// Проверяем, были ли указаны пути к файлу и директории, отличные от значений по умолчанию
	if *inputFile == " "{
		flag.PrintDefaults()
		os.Exit(1)
	}
		
	if *outputDir  == " "{
		flag.PrintDefaults()}
	return inputFile, outputDir
}


func setupOutputDir(dirPath string) {// Функция для создания директории для вывода, если она еще не существует
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("ошибка создания директории", err)// Вывод сообщения об ошибке и выход из программы при возникновении
		os.Exit(1)
	}
}

func processInputFile(filePath string, outputDir string) {// Функция для обработки входного файла с URL-адресами
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ошибка открытия файла:", err)// Вывод сообщения об ошибке и выход из программы при возникновении ошибки открытия файла
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		url := scanner.Text()
		url = strings.TrimSpace(url)
		if url != "" {
			_, err := fetchHTML(url) // Попытка получить HTML
			if err == nil { // Если нет ошибок, считаем URL валидным
				fmt.Printf("валидная ссылка: %s\n", url)
				wg.Add(1)
				go processURL(url, &wg, i, outputDir)
			} else {
				fmt.Printf("невалидная ссылка: %s\n", url)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("ошибка чтения файла:", err)// Вывод сообщения об ошибке сканирования файла и выход из программы
		os.Exit(1)
	}

	wg.Wait()
}

// Функция для обработки каждого URL в отдельной горутине
func processURL(url string, wg *sync.WaitGroup, i int, outputDir string) {
	defer wg.Done()
	fmt.Printf("Горутина для %s начата\n", url)
	html, err := fetchHTML(url)
	if err != nil {
		fmt.Printf("ошибка получения HTML для %s: %v\n", url, err)// Вывод сообщения об ошибке получения HTML и возврат из функции
		return
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("url_%d.html", i+1))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("ошибка создания файла для %s: %v\n", url, err)// Вывод сообщения об ошибке создания файла и возврат из функции
		return
	}
	defer file.Close()

	_, err = file.Write(html)
	if err != nil {
		fmt.Printf("ошибка записи HTML %s: %v\n", url, err)// Вывод сообщения об ошибке записи HTML и возврат из функции
	}
	fmt.Printf("Горутина для %s завершена\n", url)
}
