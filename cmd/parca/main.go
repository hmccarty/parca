package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	calcmd "github.com/hmccarty/parca/internal/commands/calendar"
	curcmd "github.com/hmccarty/parca/internal/commands/currency"
	vercmd "github.com/hmccarty/parca/internal/commands/verify"
	"github.com/hmccarty/parca/internal/events"
	"github.com/hmccarty/parca/internal/models"
	"github.com/hmccarty/parca/internal/services/calendar"
	"github.com/hmccarty/parca/internal/services/config"
	"github.com/hmccarty/parca/internal/services/discord"
	"github.com/hmccarty/parca/internal/services/redis"
	"github.com/hmccarty/parca/internal/services/smtp"
)

func main() {
	conf, err := config.NewConfig("config/main.yml")
	if err != nil {
		fmt.Println(err)
	}

	calendarClient := calendar.NewGoogleCalendarClient(conf)
	smtpClient, err := smtp.NewSMTPClient(conf)
	if err != nil {
		fmt.Println(err)
	}

	createDbClient := func() models.DbClient {
		return redis.OpenRedisClient(conf)
	}

	var commandList = []models.Command{
		// Currency Commands
		curcmd.NewBalanceCommand(createDbClient),
		curcmd.NewSetBalanceCommand(createDbClient),
		curcmd.NewLeaderboardCommand(createDbClient),
		curcmd.NewThanksCommand(createDbClient),
		curcmd.NewPayCommand(createDbClient),

		// Calendar Commands
		calcmd.NewAddCalendarCommand(createDbClient, calendarClient),
		calcmd.NewPrintCalendarCommand(createDbClient, calendarClient),
		calcmd.NewRemoveCalendarCommand(createDbClient, calendarClient),
		calcmd.NewTodayCommand(createDbClient, calendarClient),
		calcmd.NewWeekCommand(createDbClient, calendarClient),

		// Verification Commands
		vercmd.NewConfigureVerifyCommand(createDbClient),
		vercmd.NewVerifyCommand(createDbClient, smtpClient),
	}

	var eventList = []models.Event{
		events.NewVerifyOnMessageEvent(createDbClient),
	}

	session, err := discord.NewDiscordSession(conf, commandList, eventList)
	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}
