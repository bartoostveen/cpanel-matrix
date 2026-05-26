package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bartoostveen.nl/cpanel-matrix/config"
	"bartoostveen.nl/cpanel-matrix/matrix"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type WebhookRequest struct {
	Token    []string
	Hostname []string
	Subject  []string
	Body     []string
}

func Run(config config.AppConfig) {
	log.Info("starting application")

	matrix.InitMatrix(config.Matrix)

	mux := http.NewServeMux()
	mux.HandleFunc("/hook/{room}", createHookHandler(config))

	log.Infof("server running on :%d", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), mux))
}

func createHookHandler(config config.AppConfig) func(http.ResponseWriter, *http.Request) {
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

		token, exists := config.Matrix.Rooms[room]
		if !exists {
			http.Error(writer, "room not found", 404)
			return
		}

		var parsedRequest WebhookRequest
		err := mapstructure.Decode(request.Form, &parsedRequest)
		if err != nil {
			handleBadRequestError(writer, err, "failed to parse webhook request")
			return
		}

		if len(parsedRequest.Token) < 1 || subtle.ConstantTimeCompare([]byte(parsedRequest.Token[0]), []byte(token)) == 0 {
			http.Error(writer, "unauthorized", 401)
			return
		}

		if len(parsedRequest.Subject) != len(parsedRequest.Body) || len(parsedRequest.Subject) != len(parsedRequest.Hostname) {
			http.Error(writer, "invalid form data: amount of occurrences of distinct fields must match", 400)
			return
		}

		for i := 0; i < len(parsedRequest.Subject); i++ {
			matrix.SendMatrixMessage(room, parsedRequest.Subject[i], parsedRequest.Hostname[i], parsedRequest.Body[i])
		}
	}
}

func handleBadRequestError(writer http.ResponseWriter, err error, msg string) {
	log.WithError(err).Error(msg)
	http.Error(writer, fmt.Sprintf("Bad Request: %s", msg), 400)
}
