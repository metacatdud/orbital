package orbital

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orbital/pkg/logger"
	"path"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Route struct {
	ServiceName string
	ActionName  string
	Handler     http.HandlerFunc
	Method      string
}

type HTTPService interface {
	Register(route Route)
	OnError(w http.ResponseWriter, r *http.Request, err error)
	Use(mw ...Middleware)
}

type Server struct {
	log         *logger.Logger
	routes      map[string]Route
	notFound    http.HandlerFunc
	onError     func(w http.ResponseWriter, r *http.Request, err error)
	wsConn      *WsConn
	middlewares []Middleware
}

func NewServer(log *logger.Logger) *Server {

	srv := &Server{
		log:         log,
		routes:      make(map[string]Route),
		notFound:    onNotFoundHandler,
		onError:     onErrorHandler,
		middlewares: []Middleware{},
	}

	srv.Use(
		LoggerMiddleware(log),
		PanicRecoverMiddleware(),
	)

	return srv
}

func (s *Server) Use(mw ...Middleware) {
	s.middlewares = append(s.middlewares, mw...)
}

func (s *Server) Register(route Route) {
	routePath := path.Clean(fmt.Sprintf("/rpc/%s/%s", route.ServiceName, route.ActionName))
	s.log.Info("Register", "path", routePath)
	if _, found := s.routes[routePath]; found {
		s.log.Error("route not found", "route", route)
		return
	}

	s.routes[routePath] = route
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get Route
	route, ok := s.routes[r.URL.Path]
	if !ok {
		s.notFound.ServeHTTP(w, r)
		return
	}

	handler := route.Handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}

	handler.ServeHTTP(w, r)
}

func (s *Server) OnError(w http.ResponseWriter, r *http.Request, err error) {
	s.onError(w, r, err)
}

// Decode incoming http request body
func Decode(r io.ReadCloser, data any) error {
	bodyBytes, err := io.ReadAll(io.LimitReader(r, 1024*1024))
	if err != nil {
		return fmt.Errorf("%w", ErrBadPayload)
	}

	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		return fmt.Errorf("%w", ErrUnmarshalPayload)
	}
	return nil
}

// Encode outgoing http response
func Encode(w http.ResponseWriter, r *http.Request, status int, data any) error {
	var writer io.Writer = w

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		writer = gzw
		defer func(gzw *gzip.Writer) {
			_ = gzw.Close()
		}(gzw)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err = writer.Write(b); err != nil {
		return err
	}

	return nil
}

func onErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	_ = Encode(w, r, http.StatusInternalServerError, Error{
		Internal,
		err,
	})
}

func onNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	_ = Encode(w, r, http.StatusInternalServerError, Error{
		NotFound,
		ErrPathNotFound,
	})
}
