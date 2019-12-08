package auth

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/muerwre/orchidgo/app"
	"github.com/muerwre/orchidgo/logger"
)

// type AuthRouter struct {
// 	Router *mux.Router
// }

// Router creates new AuthRouter
func Router(router *mux.Router, logger *logger.Logger) {
	fmt.Println(router, logger)

	router.Handle("/second", gziphandler.GzipHandler(logger.Log(TestHandler)))
}

// TestHandler handles sample request
func TestHandler(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte(`its a test hand!`))
	return err
}
