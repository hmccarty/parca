package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hmccarty/parca/internal/commands/currency"
	"github.com/hmccarty/parca/internal/models"
	"github.com/hmccarty/parca/internal/services/config"
	"github.com/hmccarty/parca/internal/services/discord"
	"github.com/hmccarty/parca/internal/services/redis"
)

func main() {
	conf, err := config.NewConfig("config/main.yml")
	if err != nil {
		fmt.Println(err)
	}

	createDbClient := func() models.DbClient {
		return redis.OpenRedisClient(conf)
	}

	var commandList = []models.Command{
		currency.NewBalanceCommand(createDbClient),
		currency.NewSetBalanceCommand(createDbClient),
		currency.NewLeaderboardCommand(createDbClient),
		currency.NewThanksCommand(createDbClient),
		currency.NewPayCommand(createDbClient),
	}

	session, err := discord.NewDiscordSession(conf, commandList)
	if err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}