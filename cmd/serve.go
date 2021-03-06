package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidmap-backend/api"
	"github.com/muerwre/orchidmap-backend/app"
)

func serveAPI(ctx context.Context, api *api.API) {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	api.Init(router.Group("/api"))

	hasCerts := len(api.Config.TlsFiles) == 2

	s := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", api.Config.Host, api.Config.Port),
		Handler:     router,
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		<-ctx.Done()

		if err := s.Shutdown(context.Background()); err != nil {
			logrus.Error(err)
		}

		close(done)
	}()

	go func() {
		api.App.DB.CleanUp(nil)

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				api.App.DB.CleanUp(&t)
			}
		}
	}()

	if hasCerts {
		logrus.Infof(fmt.Sprintf("Listening https://%s:%d", api.Config.Host, api.Config.Port))

		if err := s.ListenAndServeTLS(api.Config.TlsFiles[0], api.Config.TlsFiles[1]); err != http.ErrServerClosed {
			logrus.Error(err)
		}
	} else {
		logrus.Infof(fmt.Sprintf("Listening http://%s:%d", api.Config.Host, api.Config.Port))

		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Error(err)
		}
	}

	<-done
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the api",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New()

		if err != nil {
			return err
		}

		defer a.Close()

		api, err := api.New(a)

		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			logrus.Info("signal caught. shutting down...")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)

		go func() {
			defer wg.Done()
			defer cancel()
			serveAPI(ctx, api)
		}()

		wg.Wait()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
