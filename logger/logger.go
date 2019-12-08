package logger

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/model"
	"github.com/sirupsen/logrus"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

// Logger logs all request to console
type Logger struct {
	App    *app.App
	Config *Config
}

// CreateLogger creates new Logger
func CreateLogger(a *app.App) (logger *Logger, err error) {
	logger = &Logger{
		App: a,
	}

	logger.Config, err = InitConfig()

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Log is a subrouter, that adds logging
func (l *Logger) Log(f func(*app.Context, http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		beginTime := time.Now()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctx := l.App.NewContext().WithRemoteAddress(l.IPAddressForRequest(r, l.Config))
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

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		if err := f(ctx, w, r); err != nil {
			ctx.Logger.Error(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}

// IPAddressForRequest determines IP Address for request
func (l *Logger) IPAddressForRequest(r *http.Request, c *Config) string {
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
