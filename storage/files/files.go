package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"main.go/lib/er"
	"main.go/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string // Хранит информацию о том в какой папке все будем хранить
}

const (
	defaultPerm = 0774 //0774 дает разрешение на чтение и запись
)

var ErrNoSaved = errors.New("don't have saved pages")

func New(basePath string) Storage { //Функция, которая создает базовый путь
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = er.WrapIfErr("can't save page", err) }() // способ обработки ошибок
	//Для начала определимся куда сохраняем наш файл "Путь"
	fPath := filepath.Join(s.basePath, page.UserName) // Используем filepath который определяет систему и автоматически расставляет разделители '/ unix or \ win' (storage/files)
	//Теперь этот путь нужно создать
	if err := os.Mkdir(fPath, defaultPerm); err != nil { //Mkdir создает директорию по указанному пути, получает 2 арг пермишены
		return err
	}

	fName, err := fileName(page) //Определяем название файла
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName) //Добавляем к пути файла само имя файла

	file, err := os.Create(fPath) //Создаем файл
	if err != nil {
		return err
	}
	//Закрываем созданный файл
	defer func() { _ = file.Close() }() // _=... - явно показываем что игнорируем err
	//Cериализация нашей страницы (Приведение к формату который мы можем записать в файл и по нему установить исходную структ)
	if err := gob.NewEncoder(file).Encode(page); err != nil { // страница переведена в формат gob и записана в указанный файл
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = er.WrapIfErr("can't pick random page", err) }()

	fPath := filepath.Join(s.basePath, userName) // Получаем путь до директории

	files, err := os.ReadDir(fPath) // получаем список файлов
	if err != nil {
		return nil, err
	}
	if len(files) == 0 { // Проверяем на количество файлов
		return nil, ErrNoSaved
	}
	rand.Seed(time.Now().UnixNano()) // получение рандомного числа опираясь на сид по времени
	n := rand.Intn(len(files))       // Получили само число
	file := files[n]                 // Получаем файл по номеру который сгенерировали

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p) // Получаем название файла
	if err != nil {
		return er.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName) // Получаем путь
	if err := os.Remove(path); err != nil {                 // Удаляем
		msg := fmt.Sprintf("can't remove file %s", path)

		return er.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p) // Получаем название файла
	if err != nil {
		return false, er.Wrap("can't check if file exist", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName) // Получаем путь
	switch _, err = os.Stat(path); {                        //Проверка на ошибки
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, er.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) { //функция для декодирования файла
	f, err := os.Open(filepath) // открываем файл
	if err != nil {
		return nil, er.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }() // закрываем файл

	var p storage.Page // переменная в которую будет декодирован файл

	if err := gob.NewDecoder(f).Decode(&p); err != nil { // Декодируем с помощью gob \ ! используем &p
		return nil, er.Wrap("can't decode page", err)
	}
	return &p, err
}

func fileName(p *storage.Page) (string, error) { // Функция для получения имени файла
	return p.Hash()
}
