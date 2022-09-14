package telegram

type UpdateResponce struct { // так выглядит ответ от tg
	Ok     bool     `json:"ok"`     // если true то result !=nil
	Result []Update `json:"result"` // состоит из массива данных походящих под нашу структуру ниже
}

type Update struct { // так выглядит result внутри
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct { //Структура входящего сообщения
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
