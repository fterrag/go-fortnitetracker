# go-fortnitetracker
Fortnite Tracker API Client for Go

## Getting Started

Before using this library, you'll need to obtain an API key from [https://fortnitetracker.com/site-api](https://fortnitetracker.com/site-api).

Here's some example code taken from [example/example.go](example/example.go) that will display the total number of lifetime solo kills for Ninja:

```go
httpClient := &http.Client{}
key := "your-api-key"

tracker := fortnitetracker.NewFortniteTracker(httpClient, key)
stats, _ := tracker.GetStats("pc", "ninja")

fmt.Printf("%d total solo kills\n", stats.Stats.LifetimeSolo.Kills.ValueInt)
```

[Tracker Network API Terms of Use](https://docs.google.com/document/d/1p3C7hV1WOo4figK2CNzSG_muAuszUIJ-hzzrv2toqrE/edit)
