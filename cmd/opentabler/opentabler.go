package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/goddardcm/peloscheduler/internal/opentable"
	"github.com/goddardcm/peloscheduler/internal/twilio"
)

func main() {
	appConfig, interval := mustGetConfig()

	// Set up graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-shutdownChan:
			shutdown()
		case <-ticker.C:
			queryAndNotify(appConfig)
		}
	}
}

func queryAndNotify(appConfig config.Config) {
	fmt.Printf("Querying Opentable at %s...\n", time.Now().Format(time.RFC1123))
	availability, availabilityErr := opentable.SearchForAvailability(appConfig.OpenTable)
	if availabilityErr != nil {
		fmt.Printf("Received error fetching from Opentable: %+v\n", availabilityErr)
		return
	}

	if len(availability) == 0 {
		fmt.Println("No available reservations returned from Open Table.")
		return
	}

	if err := twilio.SendMessage(
		appConfig.Twilio,
		fmt.Sprintf(
			"Open Table reservations were found: %d",
			len(availability),
		),
	); err != nil {
		fmt.Printf("Error sending message: %+v\n", err)
	}

}

func mustGetConfig() (config.Config, time.Duration) {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config file>\n", os.Args[0])
		os.Exit(1)
	}

	appConfig, configErr := config.FromFile(os.Args[1])
	if configErr != nil {
		fmt.Printf("Error: %+v\n", configErr)
		os.Exit(1)
	}

	interval, intervalErr := time.ParseDuration(appConfig.OpenTable.QueryInterval)
	if intervalErr != nil {
		fmt.Printf("Invalid query interval: %s\n", intervalErr.Error())
		os.Exit(1)
	}

	return appConfig, interval
}

func shutdown() {
	fmt.Println("Shutting down...")
	os.Exit(0)
}
