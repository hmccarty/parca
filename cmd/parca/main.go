package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// calcmd "github.com/hmccarty/parca/internal/commands/calendar"
	curcmd "github.com/hmccarty/parca/internal/commands/currency"
	gencmd "github.com/hmccarty/parca/internal/commands/general"
	vercmd "github.com/hmccarty/parca/internal/commands/verify"

	"github.com/hmccarty/parca/internal/models"
	"github.com/hmccarty/parca/internal/services/config"
	"github.com/hmccarty/parca/internal/services/discord"
	"github.com/hmccarty/parca/internal/services/redis"
	"github.com/hmccarty/parca/internal/services/smtp"
)

func main() {
	conf, err := config.NewConfig("config/main.yml")
	if err != nil {
		log.Fatal(err)
	}

	// scheduler := cron.New()
	// scheduler.Start()

	// calendarClient := gcalendar.NewGoogleCalendarClient(conf)
	smtpClient, err := smtp.NewSMTPClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	createDbClient := func() models.DbClient {
		return redis.OpenRedisClient(conf)
	}

	var commandList = []models.Command{
		// General Commands
		gencmd.NewStatusCommand(),
		gencmd.NewRoleMenuCommand(),
		gencmd.NewPollCommand(createDbClient),
		gencmd.NewBountyCommand(conf.BountyAmt, createDbClient),
		gencmd.NewRemindCommand(),
		gencmd.NewEmbedCommand(),

		// Currency Commands
		curcmd.NewBalanceCommand(createDbClient),
		curcmd.NewSetBalanceCommand(conf.ModIDs, createDbClient),
		curcmd.NewLeaderboardCommand(createDbClient),
		curcmd.NewThanksCommand(createDbClient),
		curcmd.NewPayCommand(createDbClient),

		// Calendar Commands
		// calcmd.NewAddCalendarCommand(createDbClient, calendarClient),
		// calcmd.NewPrintCalendarCommand(createDbClient, calendarClient),
		// calcmd.NewRemoveCalendarCommand(createDbClient, calendarClient),
		// calcmd.NewTodayCommand(createDbClient, calendarClient),
		// calcmd.NewWeekCommand(createDbClient, calendarClient),

		// Verification Commands
		vercmd.NewConfigureVerifyCommand(conf.ModIDs, createDbClient),
		vercmd.NewVerifyCommand(createDbClient, smtpClient),
	}

	session, err := discord.NewDiscordClient(conf, commandList, []models.Event{})

	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// scheduler.Stop()
	session.Close()
}
