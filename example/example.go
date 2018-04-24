package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fterrag/go-fortnitetracker"
)

func main() {
	httpClient := &http.Client{}
	key := "" // Generate a Fortnite Tracker API key at https://fortnitetracker.com/site-api

	tracker := fortnitetracker.NewFortniteTracker(httpClient, key)
	stats, err := tracker.GetStats("pc", "ninja")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Stats for %s\n", stats.EpicUserHandle)
	fmt.Printf("%d total solo kills\n", stats.Stats.LifetimeSolo.Kills.ValueInt)
	fmt.Printf("%d total duo kills\n", stats.Stats.LifetimeDuo.Kills.ValueInt)
	fmt.Printf("%d total squad kills\n", stats.Stats.LifetimeSquad.Kills.ValueInt)
}
