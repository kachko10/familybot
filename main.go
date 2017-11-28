package main

import (
	"os"
	"strings"
	"fmt"
	"golang.org/x/net/websocket"
	"github.com/bluele/slack"
	"path/filepath"
)

const (
	FAMILY_PREFIX = "info:"
)

const (
	begnningStory = `Noam Kachko Born in 10.09.83 Riga Latvia.
Parents Boris and Roza electrical engineer , teacher.`
	beginningFileName = "beginning.png"
	beinningCode      = "beginning"
)

const (
	wrongFamilyOptionFileName = "noEntry.png"
	wrongFamilyOptionStroy    = "no such option "
)

const (
	familyFileName = "family.png"
	familyCode     = "family"
	familyStory    = `Wife ~5 years old child and one more in the oven.`
)
const (
	divingFileName = "diving.png"
	divingCode     = "diving"
	divingStory    = "fun fun !!!"
)

const (
	skiFileName = "ski.png"
	skiCode     = "ski"
	skiStory    = "fun fun !!!"
)

const (
	biographyFileName = "taxi.png"
	biographyCode     = "biography"
	biographyStory    = `Team leader Role, last company worked in, HP.
~9 years in the software industry in many different roles:
Developer,DevOps Team Lead , Dev Team Lead. `
)

func main() {

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		panic("no slack slackToken supplied")
	}

	ws, slackBootId := slackConnect(slackToken)
	for {
		msg, err := getMessage(ws)

		if err != nil {
			println(err)
		}

		replay(msg, slackBootId, slackToken, ws, )
	}
}

func replay(m Message, slackBootId string, slackToken string, ws *websocket.Conn) {

	if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+slackBootId+">") &&
		!strings.Contains(m.Text, "uploaded a file") {
		//stock quote should be prefixed stock:"
		if strings.Contains(m.Text, FAMILY_PREFIX) {
			go func(m Message) {
				postFamilyInfo(m, slackToken)
			}(m)
		} else {

			m.Text = "Dude i can only give info about noam"
			postMessage(ws, m)
		}

	}
}

//https://www.quandl.com/api/v3/datasets/WIKI/FB.json?column_index=4&start_date=2014-01-01&end_date=2014-12-31&collapse=monthly&transform=rdiff&api_key=some_key
func postFamilyInfo(m Message, slackToken string) {

	familyTopic := strings.SplitAfter(m.Text, FAMILY_PREFIX)[1]
	familyTopic = strings.TrimSpace(familyTopic)
	api := slack.New(slackToken)
	var uploadFilePath, initialComment string
	switch familyTopic {
	case beinningCode:
		uploadFilePath = beginningFileName
		initialComment = begnningStory
	case divingCode:
		uploadFilePath = divingFileName
		initialComment = divingStory
	case skiCode:
		uploadFilePath = skiFileName
		initialComment = skiStory
	case biographyCode:
		uploadFilePath = biographyFileName
		initialComment = biographyStory
	case familyCode:
		uploadFilePath = familyFileName
		initialComment = familyStory
	default:
		uploadFilePath = wrongFamilyOptionFileName
		initialComment = wrongFamilyOptionStroy

	}

	info, err := api.FilesUpload(&slack.FilesUploadOpt{
		Filepath:       uploadFilePath,
		Filetype:       "png",
		Filename:       filepath.Base(uploadFilePath),
		Title:          "upload family",
		Channels:       []string{m.Channel},
		InitialComment: initialComment,
	})

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Completed file upload with the ID: '%s'.", info.ID))
}
