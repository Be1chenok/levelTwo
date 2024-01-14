package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

/*
	Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.

	В рамках задания необходимо:
	Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	Реализовать middleware для логирования запросов

	Методы API:
	POST /create_event
	POST /update_event
	POST /delete_event
	GET /events_for_day
	GET /events_for_week
	GET /events_for_month

	Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
	В GET методах параметры передаются через queryString, в POST через тело запроса.
	В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."}
	в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

	В рамках задачи необходимо:
	Реализовать все методы.
	Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
	В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
	В случае остальных ошибок сервер должен возвращать HTTP 500.
	Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

var (
	ErrInvalidDate        = errors.New("invalid date")
	ErrUserIDIsNotNumber  = errors.New("user_id is not number")
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrUnknownUser        = errors.New("unknown user")
	ErrEventNotFound      = errors.New("event not found")
	ErrDuplicateEvent     = errors.New("duplicate event")
	ErrInvalidPort        = errors.New("invalid port")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrSomethingWentWrong = errors.New("oops, Something went wrong")
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

// Интерфейс события
type EventStoring interface {
	CreateEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(event Event) error
	GetEventsForDay(userID int, date time.Time) ([]Event, error)
	GetEventsForWeek(userID int, date time.Time) ([]Event, error)
	GetEventsForMonth(userID int, date time.Time) ([]Event, error)
}

// Структура события
type Event struct {
	UserID      int       `json:"user_id"`
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// Структура хранилища
type eventStore struct {
	data  map[int][]Event
	mutex sync.RWMutex
}

// Конструктор хранилища событий
func NewEventStore() EventStoring {
	return &eventStore{
		data: make(map[int][]Event),
	}
}

// Создание события
func (es *eventStore) CreateEvent(event Event) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	if es.containsEvent(event) {
		return ErrDuplicateEvent
	}

	es.data[event.UserID] = append(es.data[event.UserID], event)

	return nil
}

// Обновление события
func (es *eventStore) UpdateEvent(event Event) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	events := es.data[event.UserID]
	for i := range events {
		if events[i].ID == event.ID {
			events[i].Title = event.Title
			events[i].Description = event.Description
			events[i].Date = event.Date
			return nil
		}
	}

	return ErrEventNotFound
}

// Удаление события
func (es *eventStore) DeleteEvent(event Event) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	events := es.data[event.UserID]
	for i := range events {
		if events[i].ID == event.ID {
			es.data[event.UserID] = append(es.data[event.UserID][0:i], es.data[event.UserID][i+1:]...)
			return nil
		}
	}

	return ErrEventNotFound
}

// Получение событий за день из хранилища
func (es *eventStore) GetEventsForDay(userID int, date time.Time) ([]Event, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	var result []Event
	allUserEvents := es.data[userID]
	if allUserEvents == nil {
		return nil, ErrUnknownUser
	}

	for _, event := range allUserEvents {
		if event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() &&
			event.Date.Day() == date.Day() {
			result = append(result, event)
		}
	}

	return result, nil
}

// Получение событий за неделю из хранилища
func (es *eventStore) GetEventsForWeek(userID int, date time.Time) ([]Event, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	var result []Event
	allUserEvents := es.data[userID]
	if allUserEvents == nil {
		return nil, ErrUnknownUser
	}

	for _, event := range allUserEvents {
		difference := date.Sub(event.Date)
		if difference < 0 {
			difference = -difference
		}
		if difference <= time.Duration(7*24)*time.Hour {
			result = append(result, event)
		}
	}

	return result, nil
}

// Получение событий за месяц из хранилища
func (es *eventStore) GetEventsForMonth(userID int, date time.Time) ([]Event, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	var result []Event

	allUserEvents := es.data[userID]
	if allUserEvents == nil {
		return nil, ErrUnknownUser
	}

	for _, event := range allUserEvents {
		if event.Date.Year() == date.Year() || event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result, nil
}

// Проверяет содержится ли событие в хранилище
func (es *eventStore) containsEvent(event Event) bool {
	for _, v := range es.data[event.UserID] {
		if v.ID == event.ID {
			return true
		}
	}
	return false
}

// Структура ответа с ошибкой
type errorResponse struct {
	Message string `json:"error"`
}

// Структура ответа с результатом
type resultResponse struct {
	Message interface{} `json:"result"`
}

// Пишет ответ с ошибкой
func writeJsonErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	response := errorResponse{
		Message: err.Error(),
	}
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Пишет ответ с результатом
func writeJsonResultResponse(w http.ResponseWriter, statusCode int, result interface{}) {
	response := resultResponse{
		Message: result,
	}
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Структура хэндлера
type Handler struct {
	eventStore EventStoring
}

// Конструктор хэндлера
func NewHandler(eventStore EventStoring) *Handler {
	return &Handler{
		eventStore: eventStore,
	}
}

// Пользовательский тип ResponseWriter для перехвата статуса
type responseWriterProxy struct {
	http.ResponseWriter
	statusCode int
}

// Инициализация хэндлеров
func (h Handler) InitRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/create_event", h.CreateEvent)
	router.HandleFunc("/update_event", h.UpdateEvent)
	router.HandleFunc("/delete_event", h.DeleteEvent)
	router.HandleFunc("/events_for_day", h.EventsForDay)
	router.HandleFunc("/events_for_week", h.EventsForWeek)
	router.HandleFunc("/events_for_month", h.EventsForMonth)

	loggedMux := logRequests(router)

	return loggedMux
}

// Middleware для логирования запросов
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем информацию о запросе
		log.Printf("method: %s, path: %s", r.Method, r.URL.Path)

		proxyWriter := &responseWriterProxy{ResponseWriter: w, statusCode: http.StatusOK}

		// Передаем управление следующему обработчику
		next.ServeHTTP(proxyWriter, r)

		log.Printf("status: %d", proxyWriter.statusCode)
	})
}

// Переопределение WriteHeader для перехвата статуса
func (p *responseWriterProxy) WriteHeader(code int) {
	p.statusCode = code
	p.ResponseWriter.WriteHeader(code)
}

// Хэндлер создания события
func (h Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	event, err := parseJson(r)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if eventIsValid(event) {
		if err := h.eventStore.CreateEvent(event); err != nil {
			writeJsonErrorResponse(w, http.StatusInternalServerError, ErrDuplicateEvent)
			return
		}
		writeJsonResultResponse(w, http.StatusOK, event)
		return
	}

	writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody)
}

// Хэндлер обновления события
func (h Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	event, err := parseJson(r)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	if eventIsValid(event) {
		if err := h.eventStore.UpdateEvent(event); err != nil {
			writeJsonErrorResponse(w, http.StatusInternalServerError, ErrEventNotFound)
			return
		}
		writeJsonResultResponse(w, http.StatusNoContent, nil)
		return
	}

	writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody)
}

// Хэндлер удаления события
func (h Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	if err := h.eventStore.DeleteEvent(event); err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrEventNotFound)
		return
	}

	writeJsonResultResponse(w, http.StatusNoContent, nil)
}

// Хэндлер получения событий за день
func (h Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrUserIDIsNotNumber)
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidDate)
		return
	}

	events, err := h.eventStore.GetEventsForDay(userID, date)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrUnknownUser)
		return
	}

	writeJsonResultResponse(w, http.StatusOK, events)
}

// Хэндлер получения событий за неделю
func (h Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrUserIDIsNotNumber)
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidDate)
		return
	}

	events, err := h.eventStore.GetEventsForWeek(userID, date)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrUnknownUser)
		return
	}

	writeJsonResultResponse(w, http.StatusOK, events)
}

// Хэндлер получения событий за месяц
func (h Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonErrorResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrUserIDIsNotNumber)
		return
	}

	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		writeJsonErrorResponse(w, http.StatusBadRequest, ErrInvalidDate)
		return
	}

	events, err := h.eventStore.GetEventsForMonth(userID, date)
	if err != nil {
		writeJsonErrorResponse(w, http.StatusInternalServerError, ErrUnknownUser)
		return
	}

	writeJsonResultResponse(w, http.StatusOK, events)
}

// Проверяет событие на правильность
func eventIsValid(event Event) bool {
	if event.ID <= 0 || event.UserID <= 0 || event.Title == "" || event.Description == "" {
		return false
	}

	return true
}

// Парсит json
func parseJson(r *http.Request) (Event, error) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return event, err
	}

	return event, nil
}

// Структура сервера
type Server struct {
	httpServer http.Server
}

// Конструктор сервера
func NewServer(conf *Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Handler: handler,
		},
	}
}

// Запуск сервера
func (srv *Server) Start() error {
	if err := srv.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Закрытие сервера
func (srv *Server) Shutdown(ctx context.Context) error {
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Структура конфига
type Config struct {
	Host string
	Port int
}

// Инициализация конфига
func InitConfig() (*Config, error) {
	file, err := os.Open("config.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	config := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.SplitN(row, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	portStr := config["PORT"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, ErrInvalidPort
	}
	return &Config{
		Host: config["HOST"],
		Port: port,
	}, nil
}

func main() {
	// инициализируем конфиг
	conf, err := InitConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	// Создаем хранилище
	eventStore := NewEventStore()
	// Создаем хэндлер
	handler := NewHandler(eventStore)
	// Создаем сервер
	server := NewServer(conf, handler.InitRoutes())

	go func() {
		// Запускаем сервер
		if err := server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("shutting down")

	// Закрываем сервер
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}
