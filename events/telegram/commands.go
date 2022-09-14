package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tgBot/lib/er"
	"tgBot/storage"
)

const ( //типы команд
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

//В зависимости от сообщения нужно сделать те или иные действия
//
// doCmd - Что то по типу API роутера
func (p Processor) doCmd(text string, chatID int, username string) error {
	// chatId - куда отправлять сообщения, username кто обращается
	text = strings.TrimSpace(text) // обработка сообщения удалив лишние пробелы

	log.Printf("got new command '%s' from '%s'", text, username) //Логи

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessages(chatID, msgUnknownCmd)
	}
}

func (p Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() { err = er.WrapIfErr("can't do command save page", err) }()

	page := &storage.Page{ //Подготавливаем страницу которую будем сохранять
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page) //Проверяем существует ли такая \ Если да то err
	if err != nil {
		return err
	}
	if isExists {
		p.tg.SendMessages(chatID, msgAlreadyExist) // Отправляем сообщение об ошибки дубликата
	}

	if err := p.storage.Save(page); err != nil { // Сохраняем
		return err
	}
	if err := p.tg.SendMessages(chatID, msgSaved); err != nil { // Отправляем сообщение об успешном сохранении
		return err
	}

	return nil
}

// sendRandom - отправляет рандомную ссылку пользователю
func (p Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = er.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSaved) {
		return err
	}
	if errors.Is(err, storage.ErrNoSaved) {
		return p.tg.SendMessages(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessages(chatID, page.URL); err != nil { // Отправляет сообщение с ссылкой
		return err
	}

	return p.storage.Remove(page) // Удаляет ссылку которую отправил
}

// sendHelp - отправляет справку
func (p Processor) sendHelp(chatID int) error {
	return p.tg.SendMessages(chatID, msgHelp)
}

// sendHello - отправляет приветствие
func (p Processor) sendHello(chatID int) error {
	return p.tg.SendMessages(chatID, msgHello)
}

// isAddCmd - Проверка добавленного текста
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL - проверка является ли текс url
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
