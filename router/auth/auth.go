package auth

import (
	"encoding/json"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/logger"
)

// Router creates new AuthRouter
func Router(router *mux.Router, logger *logger.Logger) {
	router.Handle("/second", gziphandler.GzipHandler(logger.Log(TestHandler)))
}

// TestHandler handles sample request
func TestHandler(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	res, _ := json.Marshal(map[string]string{
		"test": "Hello! Its working?",
	})

	http.Error(w, string(res), 401)

	return nil
}

/*
	// err := json.NewEncoder(w).Encode(map[string]string{
	// 	"test": "Hello! Its working?",
	// })

	// return err
*/
