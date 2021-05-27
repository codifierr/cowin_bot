package main

import (
	"context"
	"cowin"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"telegram"
	"time"
)

func main() {
	var (
		pinCodes               string
		telToken               string
		chatId                 string
		interval               int
		minAgeLimit            int
		minAvailableCapacity   int
		availableCapacityDose1 int
		availableCapacityDose2 int
	)
	flag.StringVar(&pinCodes, "pincodes", "", "")
	flag.StringVar(&telToken, "tel_token", "", "")
	flag.StringVar(&chatId, "chat_id", "", "")
	flag.IntVar(&interval, "interval", 15, "")
	flag.IntVar(&minAgeLimit, "min_age_limit", 18, "")
	flag.IntVar(&minAvailableCapacity, "min_available_capacity", 4, "")
	flag.IntVar(&availableCapacityDose1, "min_available_capacity_dose1", 2, "")
	flag.IntVar(&availableCapacityDose2, "min_available_capacity_dose2", 1, "")

	flag.Parse()
	if pinCodes == "" {
		log.Fatal("pincodes are mandatory")
	}
	if telToken == "" {
		log.Fatal("telegram bot token is mandatory")
	}
	if chatId == "" {
		log.Fatal("telegram chat_id is mandatory")
	}
	//print all args
	log.Println(os.Args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting cowin vaccination locater")

	date := time.Now().Format("02-01-2006")
	producer := telegram.NewProducer(chatId, telToken)
	locator := cowin.NewVacSlotLocator(ctx, date, strings.Split(pinCodes, ","), minAgeLimit, minAvailableCapacity, availableCapacityDose1, availableCapacityDose2)
	go locator.Locate(interval)
	go func() {
		for {
			r := <-locator.ResultChan
			m, err := producer.SendMessage(fmt.Sprintf("%s,%d,%s,%d,%s,%d,%d,%d", url.QueryEscape(r.CenterName), r.Pincode, r.Date, r.MinAgeLimit, r.Vaccine, r.AvailableCapacity, r.AvailableCapacityDose1, r.AvailableCapacityDose2))
			log.Println(m, err)
		}
	}()
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Received an interrupt, stopping locator...")
}