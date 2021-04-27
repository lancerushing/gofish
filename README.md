# a "go fish" card game written in go

Toy project to use the deckofcards.com API.

The deckofcards.com client throttles itself, to limit the amount of 500s it sees.  Additionally, the client attempt 5 retries after a failure.

### Run cli simulator

The program in ./cli will automatically play a "go fish" game.

```shell
$ go run ./cli/ new Brian,Nicholas,Lance
shuffling
dealing cards
Brian [KS 7C KC 7D 2H]
Lance [5D 3S AC JH KD]
Nicholas [QH JC 8H AH 6D]
deck id = 46dv225qovge
Brian [KS 7C KC 7D 2H]
Match Found: KING
Match Found: 7
Brian [2H]
Lance [5D 3S AC JH KD]
Nicholas [QH JC 8H AH 6D]
Brian [2H]
Lance [5D 3S AC JH KD]
....

Winner: Brian

```

Note: the 'new' input can be substituted for a previous deck ID 


## server - implement a api server


see https://app.swaggerhub.com/apis/nrivadeneiravericred/Cards/1.0.0

### Run it

```shell
$ go run ./server/
```
then in another shell, you can test it: 

### Test with curl
```shell
## Create Game
curl -v --request POST --header "Content-Type: application/json" \
--data '{"players": ["neo", "morpheus"]}' http://localhost:8080/games

GAME_ID="5j342ordeqqe" # save the ID  

## Get hands
curl -v --request GET  --header "Content-Type: application/json" \
http://localhost:8080/games/${GAME_ID}/players/neo

curl -v --request GET  --header "Content-Type: application/json" \
http://localhost:8080/games/${GAME_ID}/players/morpheus


## Ask for cards
curl -v --request POST  --header "Content-Type: application/json" \
--data '{"player": "morpheus", "rank":"ace"}' http://localhost:8080/games/${GAME_ID}/players/neo/fish

curl -v --request POST  --header "Content-Type: application/json" \
--data '{"player": "morpheus", "rank":"king"}' http://localhost:8080/games/${GAME_ID}/players/neo/fish

curl -v --request POST  --header "Content-Type: application/json" \
--data '{"player": "morpheus", "rank":"four"}' http://localhost:8080/games/${GAME_ID}/players/neo/fish
```

