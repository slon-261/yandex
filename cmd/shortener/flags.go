package main

import (
	"flag"
	"slon-261/yandex/cmd/config"
)

// Адрес и порт для запуска сервера
var flagRunAddr string

// Базовый адрес результирующего сокращённого URL
var flagBaseAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {

	config := config.NewConfig()

	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением http://localhost:8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", config.RunAddr, "address and port to run server")
	// регистрируем переменную flagRunAddr
	// как аргумент -b со значением http://localhost:8000 по умолчанию
	flag.StringVar(&flagBaseAddr, "b", config.BaseAddr, "address and port for base link")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
