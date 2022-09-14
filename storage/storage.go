package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgBot/lib/er"
)

type Storage interface {
	Save(p *Page) error                        //Сохранить
	PickRandom(userName string) (*Page, error) //Выбрать
	Remove(p *Page) error                      //Удалить
	IsExists(p *Page) (bool, error)            //Проверяет наличие страницы
}

var ErrNoSaved = errors.New("don't have saved pages")

type Page struct { // Основной тип данных с которым будет работать Storage
	URL      string
	UserName string
}

// Hash - уникализирует имена файлов(с помощью хеша)
func (p Page) Hash() (string, error) {
	h := sha1.New() // sha1 формат хеширования \ встроенная функция golang
	// Так как разные пользователи могут хранить одинаковые ссылки будем хэшировать по username + url
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", er.Wrap("cant calculate hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", er.Wrap("cant calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
