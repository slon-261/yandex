package main

import (
	"flag"
	"slon-261/yandex/cmd/config"
)

// Адрес и порт для запуска сервера
var flagRunAddr string

// Базовый адрес результирующего сокращённого URL
var flagBaseURL string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {

	cfg := config.NewConfig()

	//Приоритет параметров сервера должен быть таким:
	//Если указана переменная окружения, то используется она.
	//Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
	//Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
	if cfg.EnvRunAddr != "" {
		flagRunAddr = cfg.EnvRunAddr
	} else {
		// регистрируем переменную flagRunAddr
		// как аргумент -a со значением http://localhost:8080 по умолчанию
		flag.StringVar(&flagRunAddr, "a", cfg.DefaultRunAddr, "address and port to run server")
	}

	if cfg.EnvBaseURL != "" {
		flagBaseURL = cfg.EnvBaseURL
	} else {
		// регистрируем переменную flagRunAddr
		// как аргумент -b со значением http://localhost:8000 по умолчанию
		flag.StringVar(&flagBaseURL, "b", cfg.DefaultBaseURL, "address and port for base link")
	}

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
