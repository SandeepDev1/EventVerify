package main

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"strings"
	"time"
)

type Config struct {
	Token       string `json:"token"`
	MongoDBUrl  string `json:"dbUrl"`
	MongoDBName string `json:"dbName"`
}

var (
	token  string
	url    string
	dbname string
)

func main() {
	err, config := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	token = config.Token
	url = config.MongoDBUrl
	dbname = config.MongoDBName

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/addpremium", func(c telebot.Context) error {
		temp := strings.Split(c.Message().Text, " ")
		if len(temp) != 2 {
			return c.Send("Invalid format, it should be /addpremium <user>")
		}

		uid := temp[1]
		user, err := GetUser(uid)
		if err != nil {
			fmt.Println(err)
			return c.Send("user id not found please recheck and try again")
		}

		err = AddPremiumUser(uid)
		if err != nil {
			fmt.Println(err)
			return c.Send("something went wrong while adding premium")
		}

		return c.Send(fmt.Sprintf("Successfully Added premium to the user `%v`", user.Name), &telebot.SendOptions{ParseMode: "Markdown"})
	})

	b.Handle("/approve", func(c telebot.Context) error {
		temp := strings.Split(c.Text(), " ")
		if len(temp) != 3 {
			return c.Send("Invalid format it should be  /approve <event-id> <txn-id>")
		}

		uid := temp[1]
		eventData, err := GetEventDataFromId(uid)
		if err != nil {
			return err
		}

		fmt.Println(eventData.ID)
		fmt.Println(eventData.Event[0].ID)

		var teamIds []string

		for _, event := range eventData.Event {
			eventDetails := GetEventDetails(event.ID)
			if len(event.Users) != eventDetails.Members {
				continue
			}

			teamId := getUUID()
			teamIds = append(teamIds, teamId)

			eventTeam := EventDataWithTeam{
				TeamId:    teamId,
				EventID:   eventDetails.ID,
				EventName: eventDetails.Name,
				Users:     event.Users,
				TxnId:     temp[2],
			}

			err = AddEventTeam(eventTeam)
			if err != nil {
				log.Println(err)
				fmt.Println(eventTeam)
				continue
			}

			err = UpdateTeamUsers(teamId, event.Users)
			if err != nil {
				log.Println(err)
				fmt.Println(teamId)
				continue
			}
		}

		return c.Send("Successfully Approved the transaction with the teamIds being added to their profiles\n\n`"+strings.Join(teamIds, "\n")+"`", &telebot.SendOptions{ParseMode: "Markdown"})
	})

	err, _ = Connect(url)
	if err != nil {
		log.Fatal(err)
	}

	b.Start()

}
