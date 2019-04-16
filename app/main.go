package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("862051556:AAFhT0mzNWKoqLSwK4qAMQMacUxKxdjAGkw")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	apiKey := "AIzaSyAPcN6UiEstrzc_O8qv2oFVQuNOa6yhcYo"
	client, _ := maps.NewClient(maps.WithAPIKey(apiKey))
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Location != nil {
			log.Println(update.Message.Location.Latitude, update.Message.Location.Longitude)
			location := maps.LatLng{
				Lat: update.Message.Location.Latitude,
				Lng: update.Message.Location.Longitude,
			}
			r := &maps.NearbySearchRequest{
				Radius:   10000,
				Language: "ru",
				Keyword:  "cafe",
				Location: &location,
			}
			resp, _ := client.NearbySearch(context.Background(), r)
			pretty.Println(resp)
			var messageText string
			for i := 0; i < len(resp.Results); i++ {
				place := resp.Results[i]
				messageText += fmt.Sprintf(
					"%s (%s)\n%s\n%s\n\n",
					place.Name, place.Vicinity,
					strings.Join(place.Types, ", "),
					getEmojiRating(place.Rating),
				)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			bot.Send(msg)
		}
	}
}

func getEmojiRating(rate float32) string {
	rounded := int64(rate)
	switch rounded {
	case 1:
		return "⭐"
	case 2:
		return "⭐⭐"
	case 3:
		return "⭐⭐⭐"
	case 4:
		return "⭐⭐⭐⭐"
	case 5:
		return "⭐⭐⭐⭐⭐"
	default:
		return ""
	}
}

func parseLocation(location string, r *maps.NearbySearchRequest) {
	if location != "" {
		l, _ := maps.ParseLatLng(location)
		r.Location = &l
	}
}

func parsePriceLevel(priceLevel string) maps.PriceLevel {
	switch priceLevel {
	case "0":
		return maps.PriceLevelFree
	case "1":
		return maps.PriceLevelInexpensive
	case "2":
		return maps.PriceLevelModerate
	case "3":
		return maps.PriceLevelExpensive
	case "4":
		return maps.PriceLevelVeryExpensive
	}
	return maps.PriceLevelFree
}

func parsePriceLevels(minPrice string, maxPrice string, r *maps.NearbySearchRequest) {
	if minPrice != "" {
		r.MinPrice = parsePriceLevel(minPrice)
	}

	if maxPrice != "" {
		r.MaxPrice = parsePriceLevel(minPrice)
	}
}

func parseRankBy(rankBy string, r *maps.NearbySearchRequest) {
	switch rankBy {
	case "prominence":
		r.RankBy = maps.RankByProminence
		return
	case "distance":
		r.RankBy = maps.RankByDistance
		return
	case "":
		return
	}
}

func parsePlaceType(placeType string, r *maps.NearbySearchRequest) {
	if placeType != "" {
		t, _ := maps.ParsePlaceType(placeType)

		r.Type = t
	}
}
