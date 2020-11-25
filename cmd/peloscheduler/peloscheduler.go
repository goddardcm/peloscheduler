package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/goddardcm/peloscheduler/internal/peloton"
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
	fmt.Println("Querying Peloton...")
	order, orderErr := peloton.FetchOrder(appConfig.Peloton.OrderID)
	if orderErr != nil {
		fmt.Printf("Received error fetching from Peloton: %+v\n", orderErr)
		return
	}

	if len(order.AvailableDeliveries) == 0 {
		fmt.Println("No available deliveries returned from Peloton.")
		return
	}

	currentDelivery, currentDeliveryErr := order.CurrentDelivery.GetStart()
	if currentDeliveryErr != nil {
		fmt.Printf("Error getting current delivery time: %+v\n", currentDeliveryErr)
		return
	}

	firstAvailable, firstAvailableErr := order.AvailableDeliveries[0].GetStart()
	if firstAvailableErr != nil {
		fmt.Printf("Error getting available delivery time: %+v\n", firstAvailableErr)
		return
	}

	if firstAvailable.Before(currentDelivery) {
		fmt.Println("Found earlier delivery - sending notification!")
		if err := twilio.SendMessage(
			appConfig.Twilio,
			fmt.Sprintf(
				"A sooner available delivery was found: %s",
				firstAvailable.Format(time.RFC1123),
			),
		); err != nil {
			fmt.Printf("Error sending message: %+v\n", err)
		}
	} else {
		fmt.Printf("Earliest available is %s\n", firstAvailable.Format(time.RFC1123))
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

	interval, intervalErr := time.ParseDuration(appConfig.Peloton.QueryInterval)
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
