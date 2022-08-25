package server

import (
	"encoding/json"
	"httpserver/calendar"
	"httpserver/internal/server/middleware"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	apiCreateEvent    = "/create_event"
	apiUpdateEvent    = "/update_event"
	apiDeleteEvent    = "/delete_event"
	apiEventsForDay   = "/events_for_day"
	apiEventsForWeek  = "/events_for_week"
	apiEventsForMonth = "/events_for_month"
)

//Server represents server
type Server struct {
	config    *Config
	router    *http.ServeMux
	mw        *middleware.Middleware
	eventRepo calendar.EventRepository
	logger    *logrus.Logger
}

//New creates new server
func New(config *Config, eventRepo calendar.EventRepository, logger *logrus.Logger) *Server {
	s := &Server{
		config:    config,
		router:    http.NewServeMux(),
		mw:        middleware.NewMiddleware(logger, true),
		eventRepo: eventRepo,
		logger:    logger,
	}
	s.configureRouter()
	return s
}

//Run starts the server
func (s *Server) Run() error {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	s.logger.Info("Launching server on ", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		return err
	}
	return nil
}

func (s *Server) configureRouter() {
	// POST
	s.router.HandleFunc(apiCreateEvent, s.mw.Log(s.handlerCreateEvent()))
	s.router.HandleFunc(apiUpdateEvent, s.mw.Log(s.handlerUpdateEvent()))
	s.router.HandleFunc(apiDeleteEvent, s.mw.Log(s.handlerDeleteEvent()))

	// GET
	s.router.HandleFunc(apiEventsForDay, s.mw.Log(s.handlerGetEvents(calendar.EventDAY)))
	s.router.HandleFunc(apiEventsForWeek, s.mw.Log(s.handlerGetEvents(calendar.EventWEEK)))
	s.router.HandleFunc(apiEventsForMonth, s.mw.Log(s.handlerGetEvents(calendar.EventMONTH)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//handlerCreateEvent ...
func (s *Server) handlerCreateEvent() http.HandlerFunc {
	allowedMethods := []string{http.MethodPost}
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := s.validateRequestAndReadBody(w, r, allowedMethods)
		if err != nil {
			return
		}

		userID, event, err := parseCreateEventRequest(body)
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		ev, err := s.eventRepo.CreateEvent(userID, *event)
		if err != nil {
			s.respondError(w, r, internalServerErrorCode, err)
			return
		}

		s.respondResult(w, r, http.StatusCreated, ev)
	}
}

//handlerUpdateEvent updates event
func (s *Server) handlerUpdateEvent() http.HandlerFunc {
	allowedMethods := []string{http.MethodPost}
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := s.validateRequestAndReadBody(w, r, allowedMethods)
		if err != nil {
			return
		}

		event, err := parseUpdateEventRequest(body)
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		ev, err := s.eventRepo.UpdateEvent(*event)
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		s.respondResult(w, r, http.StatusOK, ev)
	}
}

//handlerDeleteEvent deletes event
func (s *Server) handlerDeleteEvent() http.HandlerFunc {
	allowedMethods := []string{http.MethodPost}
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := s.validateRequestAndReadBody(w, r, allowedMethods)
		if err != nil {
			return
		}

		data, err := parseDeleteEventRequest(body)
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		if err := s.eventRepo.DeleteEvent(data.userID, data.eventUID); err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		s.respondResult(w, r, http.StatusNoContent, "event removed")
	}
}

//handlerGetEvents returns events for date range provided
func (s *Server) handlerGetEvents(dateRange calendar.EventDateRange) http.HandlerFunc {
	allowedMethods := []string{http.MethodGet}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.validateRequsetMethod(w, r, allowedMethods); err != nil {
			return
		}

		params, err := parseGetEventQuery(r.URL.Query())
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		ev, err := s.eventRepo.GetEvents(params.userID, params.date, dateRange)
		if err != nil {
			s.respondError(w, r, badRequestErrorCode, err)
			return
		}

		s.respondResult(w, r, http.StatusOK, ev)
	}
}

//validateRequestAndReadBody checks if request method is allowed and reads request body
//Error respond occurs if there are any errors
func (s *Server) validateRequestAndReadBody(w http.ResponseWriter, r *http.Request, methods []string) ([]byte, error) {
	if err := s.validateRequsetMethod(w, r, methods); err != nil {
		s.respondError(w, r, internalServerErrorCode, err)
		return nil, err
	}
	body, err := s.readRequestBody(w, r)
	if err != nil {
		s.respondError(w, r, internalServerErrorCode, err)
		return nil, err
	}
	return body, nil
}

func (s *Server) readRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//validateRequsetMethod returns error if request method is not allowed.
func (s *Server) validateRequsetMethod(w http.ResponseWriter, r *http.Request, methods []string) error {
	m := r.Method
	for _, method := range methods {
		if m == method {
			return nil
		}
	}
	return errMethodNotAllowed
}

func (s *Server) respondResult(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	s.respond(w, r, code, map[string]interface{}{"result": data})
}

func (s *Server) respondError(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
