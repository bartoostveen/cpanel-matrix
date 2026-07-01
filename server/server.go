package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"go.bartoostveen.nl/cpanel-matrix/config"
	"go.bartoostveen.nl/cpanel-matrix/matrix"
	"golang.org/x/exp/slices"
)

type WebhookRequest struct {
	Token    []string
	Hostname []string
	Subject  []string
	Body     []string
}

func Run(config config.AppConfig) {
	log.Info("starting application")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	matrix.InitMatrix(config.Matrix)

	mux := http.NewServeMux()
	mux.HandleFunc("/hook/{room}", createHookHandler(config))
	mux.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir(config.LogsDir+"/"))))

	server := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: mux,
	}

	log.Infof("server running on :%d", config.Port)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
	})
	wg.Go(func() {
		err := server.ListenAndServe()
		log.WithError(err).Fatal("failed to listen")
	})
	wg.Go(func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanup := func() {
			log.Info("removing 7 day old files...")
			cutoff := time.Now().Add(-7 * 24 * time.Hour)

			err := filepath.WalkDir(config.LogsDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}

				info, err := d.Info()
				if err != nil {
					log.WithError(err).Warnf("failed to stat %s", path)
					return nil
				}

				if info.ModTime().Before(cutoff) {
					if err := os.Remove(path); err != nil {
						log.WithError(err).Warnf("failed to remove %s", path)
					} else {
						log.Infof("deleted old log file %s", path)
					}
				}

				return nil
			})

			if err != nil {
				log.WithError(err).Warn("log cleanup failed")
			}
			log.Info("done cleaning up, sleeping for 7 days...")
		}

		cleanup()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cleanup()
			}
		}
	})
	wg.Wait()
}

func createHookHandler(appConfig config.AppConfig) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, "invalid method", 405)
			return
		}

		room := strings.TrimSpace(request.PathValue("room"))
		if room == "" {
			http.Error(writer, "you must specify a room", 400)
			return
		}

		otherRoomID := slices.IndexFunc(appConfig.Matrix.Rooms, func(otherRoom config.MatrixRoom) bool {
			return otherRoom.ID == room
		})
		if otherRoomID == -1 {
			http.Error(writer, "room not found", 404)
			return
		}
		token := appConfig.Matrix.Rooms[otherRoomID].Token

		err := request.ParseForm()
		if err != nil {
			handleBadRequestError(writer, err, "could not parse form")
			return
		}

		var parsedRequest WebhookRequest
		err = mapstructure.Decode(request.Form, &parsedRequest)
		if err != nil {
			handleBadRequestError(writer, err, "failed to parse webhook request")
			return
		}

		if len(parsedRequest.Token) < 1 || subtle.ConstantTimeCompare([]byte(parsedRequest.Token[0]), []byte(token)) == 0 {
			http.Error(writer, "unauthorized", 401)
			return
		}

		if len(parsedRequest.Subject) < 1 || len(parsedRequest.Body) < 1 {
			http.Error(writer, "invalid form data", 400)
			return
		}

		hostname := ""
		if len(parsedRequest.Hostname) >= 1 {
			hostname = parsedRequest.Hostname[0]
		}
		matrix.SendMatrixMessage(appConfig, room, parsedRequest.Subject[0], hostname, parsedRequest.Body[0])
	}
}

func handleBadRequestError(writer http.ResponseWriter, err error, msg string) {
	log.WithError(err).Error(msg)
	http.Error(writer, fmt.Sprintf("Bad Request: %s", msg), 400)
}
