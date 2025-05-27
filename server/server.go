package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Player struct {
	playerId    string
	playerIndex int
	conn        *websocket.Conn
}

type Lobby struct {
	lobbyId    string
	players    []Player
	maxPlayers int
	host       int
}

var lobbies []Lobby

func (lobby *Lobby) addPlayer(playerId string, conn *websocket.Conn) int {
	if len(lobby.players) == lobby.maxPlayers {
		return -1
	}

	lobby.players = append(lobby.players, Player{playerId: playerId, playerIndex: len(lobby.players), conn: conn})

	return len(lobby.players) - 1
}

func (lobby *Lobby) removePlayer(playerId string) {
	var newPlayers []Player

	for _, v := range lobby.players {
		if v.playerId != playerId {
			newPlayers = append(newPlayers, v)
		}
	}

	lobby.players = newPlayers
}

func (lobby *Lobby) send(data string, playerIndexes []int) {
	for index := range playerIndexes {
		if err := lobby.players[index].conn.WriteMessage(1, []byte(data)); err != nil {
			log.Println(err)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		if messageType == 1 {
			msgData := strings.Split(strings.TrimSpace(string(p)), " ")

			log.Println(msgData)

			if msgData[0] == "host" {
				lobbyId := uuid.NewString()

				log.Println(lobbyId)

				maxPlayers, _ := strconv.Atoi(msgData[1])

				lobby := Lobby{
					lobbyId:    lobbyId,
					players:    []Player{},
					maxPlayers: maxPlayers,
					host:       0,
				}

				lobbies = append(lobbies, lobby)

				playerIndex := lobby.addPlayer(msgData[2], conn)

				lobby.send(fmt.Sprintf("join_ret %s %d", lobbyId, playerIndex), []int{playerIndex})
			} else if msgData[0] == "join" {
				lobbyId := msgData[1]
				var lobby *Lobby

				for i, v := range lobbies {
					if v.lobbyId == lobbyId {
						lobby = &lobbies[i]
					}
				}

				if lobby != nil {
					playerIndex := lobby.addPlayer(msgData[2], conn)

					lobby.send(fmt.Sprintf("join_ret %s %d", lobbyId, playerIndex), []int{playerIndex})
				}
			} else if msgData[0] == "globbs" { // TODO: Maybe implement filtering by gameId UUID
				var lobbyStr string

				for _, v := range lobbies {
					lobbyStr += v.lobbyId + ","
				}

				if err := conn.WriteMessage(1, []byte("globbs_ret "+lobbyStr)); err != nil {
					log.Println(err)
				}
			} else if msgData[0] == "invoke" { // TODO: Implement playerId specific / targeted invokes
				for index := range lobbies {
					if lobbies[index].lobbyId != msgData[1] {
						continue
					}

					for _, player := range lobbies[index].players {
						if err := player.conn.WriteMessage(1, []byte("invoke_ret "+msgData[2])); err != nil {
							log.Println(err)
						}
					}
				}
			} else if msgData[0] == "updnetvar" {
				for index := range lobbies {
					if lobbies[index].lobbyId != msgData[1] {
						continue
					}

					for _, player := range lobbies[index].players {
						if err := player.conn.WriteMessage(1, []byte("updnetvar_ret "+msgData[2]+" "+msgData[3])); err != nil {
							log.Println(err)
						}
					}
				}
			} else if msgData[0] == "leave" {]

				for index := range lobbies {
					if lobbies[index].lobbyId != msgData[1] {
						continue
					}

					lobbies[index].removePlayer(msgData[2])

					for _, player := range lobbies[index].players {
						if err := player.conn.WriteMessage(1, []byte("updnetvar_ret "+msgData[2]+" "+msgData[3])); err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Hosting server on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
