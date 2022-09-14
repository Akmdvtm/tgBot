package telegram

import (
	"errors"
	"tgBot/client/telegram"
	"tgBot/events"
	"tgBot/lib/er"
	"tgBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct { //Определяем тип Мета который будет относиться только к телеграмм
	ChatId   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor { //
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// Fetch - этот метод принимает и обрабатывает updates
func (p Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Update(p.offset, limit) //Получаем updates
	if err != nil {
		return nil, er.Wrap("can't get events", err)
	}

	if len(updates) == 0 { // Проверка на пустоту результата
		return nil, errors.New("updates count = 0 ")
	}

	res := make([]events.Event, 0, len(updates)) //аллоцируем память под результат

	for _, u := range updates { //Проходимся по всем updates и преобразовываем их в events
		res = append(res, event(u))
	}

	//Обновим внутреннее значение поля offset
	//Потому что при следующем вызове fetch мы должны получить новую порцию событий
	p.offset = updates[len(updates)-1].ID + 1
	// Так при след запросе мы получим те апдейты у которых id больше чем у последнего уже полученных
	return res, nil
}

// Process - Этот метод выполняет действия в зависимости от типа event
func (p Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message: // Работа с сообщением
		return p.processMessage(event)
	default:
		return er.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return er.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return er.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) // ok получает bool \ Если в рез не мета то false
	if !ok {
		return Meta{}, er.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

// event - функция для преобразования апдейтов в ивенты
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{ // Добавляем параметр мета
			ChatId:   upd.Message.Chat.Id,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

// fetchText - Получает текст апдейтов
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

// fetchType - Получает типы апдейтов
func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
