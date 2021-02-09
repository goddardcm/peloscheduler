package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/goddardcm/peloscheduler/internal/onemedical"
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
	fmt.Printf("Querying One Medical at %s...\n", time.Now().Format(time.RFC1123))

	apptAvailable, apptErr := onemedical.FetchAppointments(appConfig.OneMedical)
	if apptErr != nil {
		fmt.Printf("Received error fetching from One Medical: %+v\n", apptErr)
		return
	}

	if apptAvailable == "" {
		fmt.Println("No appointments available from One Medical.")
		return
	}

	fmt.Println("Found appointment - sending notification!")
	if err := twilio.SendMessage(
		appConfig.Twilio,
		fmt.Sprintf("A One Medical appointment was found on %s!", apptAvailable),
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

	interval, intervalErr := time.ParseDuration(appConfig.OneMedical.QueryInterval)
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
