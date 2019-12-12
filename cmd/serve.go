package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/gin-gonic/gin"
	"github.com/muerwre/orchidgo/api"
	"github.com/muerwre/orchidgo/app"
	"golang.org/x/crypto/acme/autocert"
)

func serveAPI(ctx context.Context, api *api.API) {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	api.Init(router.Group("/api"))

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Config.Port),
		Handler:     router,
		ReadTimeout: 2 * time.Minute,
	}

	if len(api.Config.TlsHosts) > 0 {
		fmt.Printf("We have certs! %v", api.Config.TlsHosts)

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(api.Config.TlsHosts...), //Your domain here
			Cache:      autocert.DirCache("certs"),                     //Folder for storing certificates
		}

		s.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
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

	logrus.Infof("Listening http://127.0.0.1:%d", api.Config.Port)

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Error(err)
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
