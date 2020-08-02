package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nicklaw5/helix"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// Versioning 2.0.0 spec (https://semver.org/).
const (
	Major uint = 1
	Minor uint = 0
	Patch uint = 0
)

var version = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

type config struct {
	TwitchGenToken     bool
	TwitchClientID     string
	TwitchClientSecret string
	TwitchAccessToken  string
	bind               string
}

func parseFlags() (config, error) {
	c := config{}
	flag.BoolVar(&c.TwitchGenToken, "twitch-gen-token", false, "Generate the twitch client token")
	flag.StringVar(&c.TwitchClientID, "twitch-client-id", "", "The client id for twitch api")
	flag.StringVar(&c.TwitchClientSecret, "twitch-client-secret", "", "The secret key for twitch api")
	flag.StringVar(&c.TwitchAccessToken, "twitch-client-token", "", "Client Token for twitch api")
	flag.StringVar(&c.bind, "bind", "0.0.0.0:8080", "The host:port to bind to.")
	flag.Parse()

	if c.TwitchGenToken == true {
		if c.TwitchClientID == "" {
			return c, fmt.Errorf("you must specify --twitch-client-id\n")
		}
		if c.TwitchClientSecret == "" {
			return c, fmt.Errorf("you must specify --twitch-client-secret\n")
		}
	}
	if c.TwitchClientID == "" {
		return c, fmt.Errorf("you must specify --twitch-client-id")
	}
	return c, nil
}

func main() {
	c, err := parseFlags()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	if c.TwitchGenToken == true {
		genAccessToken(c)
		os.Exit(0)
	}
	http.HandleFunc("/twitch", c.twitch)
	http.HandleFunc("/healthz", healthHandler)
	log.Printf("Starting up Stat-Tracker [Version %s, Build: %s %s/%s]\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	log.Printf("Serving on %s\n", c.bind)
	log.Fatal(http.ListenAndServe(c.bind, nil))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func (cnf config) twitch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Twitch: Failed to parse request: %v", err.Error())
		return
	}
	streams := r.Form.Get("streams")
	if streams == "" {
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	log.Printf("Twitch: Requesting stats for streams %s)n", streams)

	w.Header().Set("Content-Type", "application/json")
	streamInfo, err := getTwitchInfo(cnf, streams)

	if err != nil {
		log.Printf("Twitch: Failed to grab stream info (%v)\n", err.Error())
		return
	}
	if streamInfo == nil {
		log.Printf("Twitch: No stream found (%s)\n", streams)
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(streamInfo)
	if err != nil {
		log.Printf("Twitch: Failed to marshal json: %v.\n", err)
	}
	_, _ = w.Write(j)

}

func getTwitchInfo(c config, s string) (stream []helix.Stream, err error) {
	// Setup client with the token
	client, err := helix.NewClient(&helix.Options{
		ClientID:       c.TwitchClientID,
		AppAccessToken: c.TwitchAccessToken,
	})
	if err != nil {
		log.Printf("GetTwitchInfo: Stream client setup failed. Please check Twitch ID & access Token (%v)", err.Error())
		return []helix.Stream{}, err
	}
	streams := strings.Split(s, ",")

	resp, err := client.GetStreams(&helix.StreamsParams{
		UserLogins: streams,
	})

	if err != nil {
		log.Printf("GetTwitchInfo: Failed to grab stream info (%v)", err.Error())
		return []helix.Stream{}, err
	}

	return resp.Data.Streams, err
}

func genAccessToken(AppConfig config) (status bool) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     AppConfig.TwitchClientID,
		ClientSecret: AppConfig.TwitchClientSecret,
	})
	if err != nil {
		log.Printf("GenAccessToken: Failed to setup client (%v)", err.Error())
		return
	}

	if client == nil {
		log.Printf("GenAccessToken: Client setup is nil (bad data?)\n")
		return
	}

	resp, err := client.GetAppAccessToken()
	if err != nil {
		log.Printf("GenAccessToken: Twitch API Error: %+v\n", err.Error())
		return false
	}

	log.Printf("Generated Client Token: %+v\n", resp.Data.AccessToken)
	log.Printf("Token will expire: %+v\n", resp.Data.ExpiresIn)
	return true
}
