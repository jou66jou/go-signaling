package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"p2p/internal/config"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

type Fields map[string]interface{}

// ApiserverCmd 是 apiserver 服務的進入口
var ApiserverCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if r := recover(); r != nil {
				// unknown error
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("unknown error: %v", err)
				}
				log.Println(err)
				time.Sleep(3 * time.Second) // wait to print all logs
			}
		}()
		cfg := config.New("app.yml")
		err = initialize(cfg)
		if err != nil {
			log.Panicf("main: initialize failed: %v", err)
			return
		}
		config.SetGlobalConfiguration(cfg)
		// fix gorm NowFunc
		gorm.NowFunc = func() time.Time {
			return time.Now().UTC()
		}

		// start http server
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		httpServer := &http.Server{
			Addr: ":" + port,
			Handler: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					method := r.Header.Get("X-HTTP-Method-Override")
					if method != "" {
						r.Method = method
					}
					next.ServeHTTP(w, r)
				})
			}(newEchoHandler(cfg)),
		}

		go func() {
			// service connections
			log.Printf("main: Listening and serving HTTP on %s\n", httpServer.Addr)
			err = httpServer.ListenAndServe()
			if err != nil {
				log.Panicf("main: http server listen failed: %v\n", err)
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		log.Println("main: shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Println("main: http server shutdown error: %v", err)
		} else {
			log.Println("main: gracefully stopped")
		}
	},
}
