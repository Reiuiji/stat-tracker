# stat-tracker

`stat-tracker` is a program that will handle the API of other platforms (Ex: twitch). For tracking the stats of streams.

### Setup Twitch API
In order to use the Twitch Helix API. You will need to generate a token from your Client API & Secret.

`go run stat-tracker.go --twitch-gen-token=true --twitch-client-id=[clientID] --twitch-client-secret=[clientSecret]` 

This will generate a Client Token you will need to run all the stat commands.

### Run stat-tracker
Run the stat-tracker with the client token to properly link with the api.
`go run stat-tracker.go --twitch-client-id=[clientID] --twitch-client-token=[clientToken]`

Once the stat-tracker is up. You can access the data through the HTTP interface. Below is an example of how to input the data. The result data will output in json format that can easily be process in a web app.

http://0.0.0.0:8080/twitch?streams=relaxbeats,twitchplayspokemon

