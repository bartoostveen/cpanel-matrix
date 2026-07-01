package matrix

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"errors"
	"fmt"
	htmlTemplate "html/template"
	"os"
	"strings"
	"sync"
	"syscall"
	textTemplate "text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bartoostveen.nl/cpanel-matrix/config"
	"go.bartoostveen.nl/cpanel-matrix/util"
	"golang.org/x/exp/slices"
	"golang.org/x/net/context"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

//go:embed templates/*
var templateFS embed.FS

var (
	htmlTemplates = htmlTemplate.Must(htmlTemplate.ParseFS(templateFS, "templates/*.html"))
	textTemplates = textTemplate.Must(textTemplate.ParseFS(templateFS, "templates/*.txt"))

	client *mautrix.Client
	ctx    = context.Background()
)

func InitMatrix(config config.MatrixConfig) {
	var err error
	client, err = mautrix.NewClient(config.HomeserverURL, id.UserID(config.MxID), config.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	client.Store = NewDefaultFileSyncStore()

	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	registerEvents(syncer, config.Rooms)

	syncCtx, cancelSync := context.WithCancel(ctx)
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)

	go func() {
		defer syncStopWait.Done()

		err = client.SyncWithContext(syncCtx)
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	util.OnSignal(func() {
		log.Infoln("Stopping Matrix sync")

		cancelSync()
		syncStopWait.Wait()
	}, syscall.SIGTERM, syscall.SIGINT)
}

func registerEvents(syncer *mautrix.DefaultSyncer, rooms []config.MatrixRoom) {
	syncer.OnEventType(event.StateMember, func(ctx context.Context, evt *event.Event) {
		roomID := evt.RoomID.String()
		log.Infof("Got asked to join room %s", roomID)

		if evt.GetStateKey() != client.UserID.String() ||
			evt.Content.AsMember().Membership != event.MembershipInvite ||
			!slices.ContainsFunc(rooms, func(room config.MatrixRoom) bool {
				return room.ID == roomID
			}) {
			return
		}

		log.Infof("Joining room %s", roomID)
		_, err := client.JoinRoomByID(ctx, evt.RoomID)
		if err != nil {
			log.WithError(err).Warn("could not join room %s", evt.RoomID)
			return
		}
		log.Infof("Joined room %s", roomID)
	})
}

type RenderRequest struct {
	Subject  string
	Hostname string
	Body     string
	Url      *string
}

func renderMatrixMessage(cfg config.AppConfig, subject string, hostname string, body string) (err error, plain string, formatted string) {
	lines := strings.Split(body, "\n")
	var msg string
	var url *string = nil

	if len(lines) > 2 {
		err, u := saveMessage(cfg, body)
		if err != nil {
			return err, "", ""
		}
		url = &u
		msg = lines[0]
	} else {
		msg = body
	}

	req := RenderRequest{
		Subject:  subject,
		Hostname: hostname,
		Body:     msg,
		Url:      url,
	}

	var plainResult bytes.Buffer
	if err = textTemplates.ExecuteTemplate(&plainResult, "message.txt", req); err != nil {
		return
	}

	var formattedResult bytes.Buffer
	if err = htmlTemplates.ExecuteTemplate(&formattedResult, "message.html", req); err != nil {
		return
	}

	plain = plainResult.String()
	formatted = formattedResult.String()
	return
}

func saveMessage(cfg config.AppConfig, msg string) (err error, url string) {
	_ = os.Mkdir(cfg.LogsDir, os.FileMode(0700))
	fileName := fmt.Sprintf("%d-%x.txt", time.Now().UnixMilli(), sha256.Sum256([]byte(msg)))
	err = os.WriteFile(cfg.LogsDir+"/"+fileName, []byte(msg), os.FileMode(0700))
	if err != nil {
		return
	}

	url = cfg.BaseUrl + "/logs/" + fileName
	return
}

func SendMatrixMessage(cfg config.AppConfig, room string, subject string, hostname string, body string) {
	err, plain, formatted := renderMatrixMessage(cfg, subject, hostname, body)
	if err != nil {
		log.WithError(err).Warn("Could not render message for Matrix!")
		return
	}
	formatted = strings.ReplaceAll(strings.TrimSpace(formatted), "\n", "<br />")

	_, err = client.SendMessageEvent(ctx, id.RoomID(room), event.EventMessage, &event.MessageEventContent{
		MsgType:       event.MsgNotice,
		Body:          plain,
		FormattedBody: formatted,
	})

	if err != nil {
		log.WithError(err).Warn("Could not deliver message to Matrix!")
	}
}
