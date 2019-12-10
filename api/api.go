package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/model"
)

type API struct {
	App    *app.App
	Config *Config
}

type ErrorCode struct {
	Code   string   `json:"code"`
	Stack  []string `json:"stack"`
	Reason string   `json:"reason"`
}

func New(a *app.App) (api *API, err error) {
	api = &API{App: a}

	api.Config, err = InitConfig()

	if err != nil {
		return nil, err
	}

	return api, nil
}

func (a *API) Init(r *mux.Router) {
	r.Use(gziphandler.GzipHandler, a.loggingMiddleware)

	(&AuthRouter{}).Init(r.PathPrefix("/auth").Subrouter(), a)
}

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (a *API) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		beginTime := time.Now()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctx := a.App.NewContext().WithRemoteAddress(a.IPAddressForRequest(r, a.Config))
		ctx = ctx.WithLogger(ctx.Logger.WithField("request_id", base64.RawURLEncoding.EncodeToString(model.NewId())))

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			logger := ctx.Logger.WithFields(logrus.Fields{
				"duration":    duration,
				"status_code": statusCode,
				"remote":      ctx.RemoteAddress,
			})
			logger.Info(r.Method + " " + r.URL.RequestURI())
		}()

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// IPAddressForRequest determines IP Address for request
func (a *API) IPAddressForRequest(r *http.Request, c *Config) string {
	addr := r.RemoteAddr

	if c.ProxyCount > 0 {
		h := r.Header.Get("X-Forwarded-For")

		if h != "" {
			clients := strings.Split(h, ",")

			if c.ProxyCount > len(clients) {
				addr = clients[0]
			} else {
				addr = clients[len(clients)-c.ProxyCount]
			}
		}
	}

	return strings.Split(strings.TrimSpace(addr), ":")[0]
}

// Handler catches and reports errors
func (a *API) Handler(f func(*app.Context, http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := a.App.NewContext().WithRemoteAddress(a.IPAddressForRequest(r, a.Config))
		ctx = ctx.WithLogger(ctx.Logger.WithField("request_id", base64.RawURLEncoding.EncodeToString(model.NewId())))

		var stack string

		if a.Config.Debug {
			stack = string(debug.Stack())
		}

		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(&ErrorCode{Code: "internal", Stack: strings.Split(stack, "\n")})
			}
		}()

		if err := f(ctx, w, r); err != nil {
			ctx.Logger.Error(err)
			reason := fmt.Sprintf("%v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				&ErrorCode{
					Code:   "internal",
					Stack:  strings.Split(stack, "\n"),
					Reason: reason,
				},
			)
			return
		}
	})
}
