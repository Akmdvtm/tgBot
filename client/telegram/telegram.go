package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tgBot/lib/er"
)

type Client struct {
	host     string      // Host tg API \ !tg-bot.com
	basePath string      //базовый путь \ Это тот префикс с которого начинаются все запросы / !host(tg-bot.com)/bot<token>...
	client   http.Client // http client \ чтобы не создавать для каждого запроса отдельно
}

const (
	getUpdatesMethod = "getUpdates" // На случай если тг изменит префикс
	sendMessageMethod
)

// New - создает client
func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath - создает новый путь
func newBasePath(token string) string { //Если нужно будет создавать еще BasePath то делать будем функц или тг изменит префикс то меняем только тут
	return "bot" + token // "bot" префикс
}

// Наш клиент занимается 2 вещами
// 1. Получает update (новые сообщения)
// 2. Отправка собственных сообщений

// Update - получение новых обновлений
func (c *Client) Update(offset, limit int) (updates []Update, err error) { // GettingUpdates - getUpdates (tgAPI\site) \ их всегда нужно указывать при вызове метода
	// Параметры запроса, удобнее это делать с помощью пакета URL
	defer func() { err = er.WrapIfErr("request error", err) }() // Если исп то нужно в return знач добавить названия (updates []Update)
	q := url.Values{}                                           // Формируем параметры запроса
	q.Add("offset ", strconv.Itoa(offset))                      // Добавляем указанный параметр к запросу
	q.Add("limit", strconv.Itoa(limit))                         // тоже =)

	//Делаем запрос
	//Так как код для запросов нашего клиента будет везде одинаковый выведем его в отдельную функцию doRequest
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	} // Получили данные из ответа
	var res UpdateResponce //Результат будем хранить в этой переменной

	//Мы понимаем что там будет JSON по этому нам нужно его распарсить
	if err := json.Unmarshal(data, &res); err != nil { // data - что парсим, &res - куда !Обязательно '&' иначе не сможет ничего добавить
		return nil, err
	}

	return res.Result, nil
}

// SendMessages - отправка сообщений
func (c *Client) SendMessages(chatID int, text string) error { // аргументы такие потому что в документации так написанно
	q := url.Values{}
	q.Add("chatId", strconv.Itoa(chatID))
	q.Add("text", text)

	//Делаем запрос \ тут тело ответа не понадобится
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return er.Wrap("can't send message:", err) // Тут просто wrap потому что ошибка 100% есть
	}

	return nil
}

// doRequest - создает запрос
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = er.WrapIfErr("request error", err) }() // работает в конце функции всегда, делает врап ошибки (в обработке)
	//Сформируем url на который будет наш запрос
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method), // Что-бы не мучатся с '/' в url
	}
	// Сформируем объект запроса \ Еще не отправляем, а только подготавливаем его
	req, err := http.NewRequest(http.MethodGet, u.String(), nil) // что бы не писать каждый раз "GET" исп http.MethodGet,
	// nil потому что тело уже есть и у http.MethodGet в основном тело отсутствует
	if err != nil { // Обработка ошибки
		return nil, err
	}
	// Теперь нужно в объект req передать параметр запроса который получили в аргументе query
	req.URL.RawQuery = query.Encode() // Encode приведет параметры к виду который мы сможем потом передать на сервер

	//Теперь отправляем получившийся запрос
	//Для отправки мы используем тот клиент который заранее подготовили
	resp, err := c.client.Do(req)
	if err != nil { // Обработка ошибки
		return nil, err
	}
	//Закрывает тело
	defer func() { _ = resp.Body.Close() }() // ошибку проигнорировали
	//Получаем содержимое
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//Осталось лишь вернуть результат
	return body, nil
}
