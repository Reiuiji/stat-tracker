package main

import (
	"fmt"
	"github.com/nicklaw5/helix"
)

func test_getViewerCount(c config, s string) (cnt int, err error) {
	// Setup client with the token
	client, err := helix.NewClient(&helix.Options{
		ClientID:       c.TwitchClientID,
		AppAccessToken: c.TwitchAccessToken,
	})
	if err != nil {
		fmt.Printf("Stream client setup failed. Please check Twitch ID & access Token (%v)", err.Error())
		return 0, err
	}

	resp, err := client.GetStreams(&helix.StreamsParams{
		UserLogins: []string{s},
	})

	if err != nil {
		return 0, err
	}

	for _, stream := range resp.Data.Streams {
		fmt.Printf("Name: %s - %d viewers\n", stream.UserName, stream.ViewerCount)
		return stream.ViewerCount, err
	}
	return 0, err
}

func test_PrintstreamInfo(c config, s string) {
	streamInfo, err := getTwitchInfo(c, s)

	if err != nil {
		fmt.Printf("Failed to grab stream info (%v)", err.Error())
		return
	}

	for _, stream := range streamInfo {
		fmt.Printf("Name: %s - %d viewers\n", stream.UserName, stream.ViewerCount)
	}
}

func test_getViewerCount2(c config) {
	// Setup client with the token
	client, err := helix.NewClient(&helix.Options{
		ClientID:       c.TwitchClientID,
		AppAccessToken: c.TwitchAccessToken,
	})
	if err != nil {
		fmt.Printf("Failed : client")
		return
	}
	isValid, _, err := client.ValidateToken(c.TwitchAccessToken)
	if err != nil {
		fmt.Printf("Failed : token")
		return
	}

	if !isValid {
		fmt.Printf("%s access token is not valid!\n", c.TwitchAccessToken)
		return
	}

	resp, err := client.GetStreams(&helix.StreamsParams{
		UserLogins: []string{"relaxbeats", "twitchplayspokemon"},
	})

	if resp == nil {
		fmt.Printf("Response code is nil")
		return
	}

	if err != nil {
		fmt.Printf("Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())
	/*
		for _, user := range resp.Data.Users {
			fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
		}*/

	for _, stream := range resp.Data.Streams {
		fmt.Printf("Name: %s - %d viewers\n", stream.UserName, stream.ViewerCount)
	}
}
