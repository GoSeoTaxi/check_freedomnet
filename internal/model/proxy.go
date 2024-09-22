package model

type Proxy struct {
	Type string // Протокол, например "socks5"
	Host string // Хост прокси-сервера
	Port string // Порт прокси-сервера
	User string // Логин
	Pass string // Пароль
}
