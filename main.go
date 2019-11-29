package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"./taonode"
)

type GameState int32

const (
	WAGER    GameState = 0
	COMEOUT  GameState = 1
	CRAPS    GameState = 2
	ON       GameState = 3
	PASSWIN  GameState = 4
	PASSLOSE GameState = 5
)

type Player struct {
	Username string `json:"username,omitempty"`
	Bet      string `json:"bet,omitempty"`
	Id       string `json:"id,omitempty"`
}

type Players struct {
	Players  []Player `json:"players,omitempty"`
	Pot      int      `json:"pot,omitempty"`
	Balance  int      `json:"balance,omitempty"`
	Point    int      `json:"point,omitempty"`
	Shooter  string   `json:"shooter,omitempty"`
	Roll     bool     `json:"roll,omitempty"`
	Dice     Dice     `json:"dice,omitempty"`
	PlaceBet bool     `json:"placeBet,omitempty"`
}

type Round struct {
	Users     []string  `json:"users,omitempty"`
	GameData  Players   `json:"gameData,omitempty"`
	State     GameState `json:"state,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	Wager     int       `json:"wager,omitempty"`
	Jumbotron string    `json:"jumbotron,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type Bet struct {
	Name  string
	Id    *websocket.Conn
	Bet   string
	Wager int
	Address string
	Balance int
}

type Dice struct {
	L     int `json:"l,omitempty"`
	R     int `json:"r,omitempty"`
	Total int `json:"total,omitempty"`
}

type LastTx struct {
	Type      string `json:"type"`
	Addresses string `json:"addresses"`
}

type TaoExplorer struct {
	Address  string   `json:"address"`
	Sent     int      `json:"sent"`
	Received string   `json:"received"`
	Balance  string   `json:"balance"`
	lastTxs  []LastTx `json:"last_txs"`
}

type FmTao struct {
	Result      string `json:"result"`
}

type WalletJSON struct {
	Action string `json:"action"`
	Address string `json:"address"`
	Withdraw string `json:"withdraw"`
	Recipient string `json:"recipient"`
	Balance string `json:"balance"`
	UUID    string `json:"uuid"`
}


var wallet WalletJSON
var Man map[string]Bet
var Game Round

var UpdateUser Round
var Arr []string
var f Player
var d Round

var Playing []string
var GState GameState
var point int
var t int
	var leftDice int
	var rightDice int
	var finishL int
	var finishR int
var shooter string
var Clients map[int]string
var N int

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func RemoveIndex(s []string, index int) []string {
	copy(s[index:], s[index+1:]) // Shift a[i+1:] left one index.
	s[len(s)-1] = ""             // Erase last element (write zero value).
	s = s[:len(s)-1]             // Truncate slice.
	return s
}



func reader(conn *websocket.Conn) {

	var ok bool


	//Read message
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			for idx, element := range Man {
				if element.Id == conn {
					for i := 0; i < len(Arr); i++ {
						if Arr[i] == element.Name {
							Arr = RemoveIndex(Arr, i)
						}
					}
					for y := 0; y < len(Game.Users); y++ {
						if Game.Users[y] == element.Name {
							Game.Users = RemoveIndex(Game.Users, y)
							d.Users = RemoveIndex(d.Users, y)
						}
					}
					for z := 0; z < len(Playing); z++ {
						if Playing[z] == element.Name {
							Playing = RemoveIndex(Playing, z)
						}
					}
					for _, element := range Man {
						element.Id.WriteJSON(Game)
					}
					_, ok := Man[idx]
					if ok {
						fmt.Println(Man)
						delete(Man, idx)
						fmt.Println(Man)
					}
					for idx, e := range Clients {
						if Man[Clients[idx]].Name == Man[e].Name {
							fmt.Println(Clients)
							delete(Clients, idx)
							fmt.Println(Clients)
							N = N - 1
						}
					}
				}
			}
			conn.Close()
			return
		}
		//unmarshal json
		if err := json.Unmarshal([]byte(p), &d); err != nil {
			panic(err)
		}
        leftDice = 0
	    rightDice = 0
	    t = 0


		if(d.Message == "rolling") {

			if d.GameData.Dice.L > 0 && d.GameData.Dice.R > 0 {

			leftDice = d.GameData.Dice.L
			rightDice = d.GameData.Dice.R
			Game.GameData.Dice.L = leftDice
			Game.GameData.Dice.R = rightDice
		    }

			for _, element := range Man {
				if element.Name != Arr[0] {
				element.Id.WriteJSON(Game)
			}
			}
		}

		Game.GameData.PlaceBet = false
		if Game.GameData.Pot > 0 {
			d.GameData.Pot = Game.GameData.Pot
		}



		//add or update player
		_, ok = Man[d.UUID]
		if !ok {

			var data FmTao
			nodeAddr := taonode.GetAddress()

			res := taonode.Balance(nodeAddr)

			if err := json.Unmarshal([]byte(res), &data); err != nil {
				fmt.Println(err)
				return
			}

			bal, err := strconv.Atoi(data.Result)
			if err != nil {
				fmt.Println(err)
			}

			Man[d.UUID] = Bet{Name: d.Message, Id: conn, Address: nodeAddr, Balance: bal, }
			fmt.Println(Man[d.UUID])


			//User entered chatroom

			Arr = append(Arr, d.Message)
			Game.Users = Arr
		    d.Users = Game.Users
			Clients[N] = d.UUID
			N++
			rows := len(Game.Users)

			shooter = Clients[0]
			Game.GameData.Shooter = Man[Clients[0]].Name
			d.GameData.Shooter = Man[Clients[0]].Name

			if rows < 2 {
				for _, element := range Arr {
					// element is the element froM someSlice for where we are


					player := Player{Username: element}
					Game.GameData.Players = append(Game.GameData.Players, player)
				}

			}



		} else if ok {

			//count bets
			if d.Message == "PASS" && d.State == WAGER || d.Message == "DONTPASS" && d.State == WAGER {

				d.Jumbotron = "Wait for others to match..."


				Man, Playing, d = placeBet(Man, Playing, d)
				Game.GameData.PlaceBet = true
				d.GameData.PlaceBet = false


				if len(Playing) == len(Arr) {


					ready := countBets(Man)
					if ready {
						d.GameData.Roll = true
						d.GameData.PlaceBet = false
						d.Jumbotron = "Roll when ready"
						Game.Jumbotron = "Bet Placed"
						Game.State = COMEOUT
						d.State = COMEOUT
					}
				} else if d.UUID != shooter {
					matched := matchedBet(Man)

					if matched {
						UpdateUser.Users = Game.Users
						UpdateUser.GameData.PlaceBet = false
						UpdateUser.Jumbotron = "Wait for others to match..."
						UpdateUser.State = COMEOUT
						Man[d.UUID].Id.WriteJSON(UpdateUser)
						continue
					} else if d.Wager != Man[shooter].Wager {
						UpdateUser.Users = Game.Users
						UpdateUser.GameData.PlaceBet = true
						UpdateUser.Wager = 0
						UpdateUser.Jumbotron = "match rollers amount:" + strconv.Itoa(Man[shooter].Wager)
						UpdateUser.State = COMEOUT
						Man[d.UUID].Id.WriteJSON(UpdateUser)
						continue
					}
				}
			}
		}




		//fnishd the dice
		if d.Message == "finished" {
			finishL = d.GameData.Dice.L
			finishR = d.GameData.Dice.R
			Game.GameData.Dice.L = finishL
			Game.GameData.Dice.R = finishR
			t = finishL + finishR + 2
			Game.GameData.Dice.Total = t
			d.GameData.Dice.Total = t
		}

		if d.Message == "finished" && d.State == COMEOUT {

			point = t
			Game.GameData.Point = point
		}




		//check state and update players

		GState = d.State
		Game.UUID = ""




		if d.Message == "finished" {
			if t == point && GState == ON {
				Game.Jumbotron = "Pass bet wins!"
				d.Jumbotron = "Pass bet wins!"


				GState = PASSWIN
				Game.State = GState
				d.State = GState
				Man, Game, d, Clients, shooter = payout(Man, Game, d, GState, Clients, shooter)

			}
			if t == 7 && GState == ON {
				Game.Jumbotron = "Don't pass bet wins!"
				d.Jumbotron = "Don't pass bet wins!"


				GState = PASSLOSE
				Game.State = GState
				d.State = GState
				Man, Game, d, Clients, shooter = payout(Man, Game, d, GState, Clients, shooter)


			}
			if GState == ON && t != point && GState== ON && t != 7 {
				Game.Jumbotron = "Point on!"
				d.Jumbotron = "Point on!"
				d.GameData.Roll = true

			}
			if GState == PASSWIN || GState == PASSLOSE || GState == CRAPS {

			}





			if t > 3 && t < 11 && t != 7 && GState == COMEOUT {

				Game.Jumbotron = "Point on"
				d.Jumbotron = "Point on"

				GState = ON
				Game.State = GState
				d.State = ON

			}
			if t == 7 && GState == COMEOUT || GState == COMEOUT && t == 11 {
				Game.Jumbotron = "Pass bet wins!"
				d.Jumbotron = "Pass bet wins!"


				GState = PASSWIN
				Game.State = GState
				d.State = GState
				Man, Game, d, Clients, shooter = payout(Man, Game, d, GState, Clients, shooter)

			}
			if t == 2 && GState == COMEOUT || GState == COMEOUT && t == 3 || GState == COMEOUT && t == 12 {

				Game.Jumbotron = "CRAPS"
				d.Jumbotron = "CRAPS"




				GState = CRAPS
				Game.State = GState

				d.GameData.PlaceBet = true
				d.State = GState
				Man, Game, d, Clients, shooter = payout(Man, Game, d, GState, Clients, shooter)
				fmt.Println("Shhoootterrr:", shooter)
			}



		switch Game.State {
		case WAGER:


		case COMEOUT:
			Game.Jumbotron = "Comeout roll"

		case CRAPS:

			d.State = CRAPS
			d.GameData.Roll = false

		case ON:


		case PASSWIN:

			d.GameData.Roll = false

		case PASSLOSE:

			d.GameData.Roll = false

		}

}
if d.GameData.Pot == 0 && d.State == WAGER {

	d.Jumbotron = "Place bet..."
d.GameData.PlaceBet = true
Game.Jumbotron = "Wait for Roller to wager.."
			Game.GameData.Roll = false
			d.GameData.Roll = false
}

if d.GameData.Pot > 0 && GState == WAGER {

	Game.GameData.PlaceBet = true
	Game.Jumbotron = "Place bet..."
}


d.GameData.Point = point




if GState == CRAPS || GState == PASSWIN || GState == PASSLOSE {
	d.GameData.PlaceBet = true
	d.Wager = 0
	Game.Wager = 0
	Playing = []string{}
	GState = WAGER
	d.State = WAGER
	Game.State = WAGER
	Game.GameData.Shooter = Man[shooter].Name
	d.GameData.Shooter = Man[shooter].Name

}




		for index, element := range Man {
			if index != Clients[0] {
				Game.GameData.Balance = Man[index].Balance
				element.Id.WriteJSON(Game)
				continue
			} else if index == Clients[0] {
				d.GameData.Balance = Man[index].Balance
				element.Id.WriteJSON(d)
			}
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//log.Println("Client Successfully Connected...")

	reader(ws)

}



func payout(M map[string]Bet, Game Round, d Round, finishState GameState, Clients map[int]string, shooter string) (map[string]Bet, Round, Round, map[int]string, string) {

	var i int
	var y int
	var name string
	var lastName string


	for _, v := range M {
		if v.Bet == "PASS" && finishState == PASSWIN || v.Bet == "DONTPASS" && finishState == PASSLOSE || v.Bet == "DONTPASS" && finishState == CRAPS {
				i++
			}
	}

	for index, v := range M {

			if v.Bet == "PASS" && finishState == PASSWIN || v.Bet == "DONTPASS" && finishState == PASSLOSE || v.Bet == "DONTPASS" && finishState == CRAPS {

				v.Balance += Game.GameData.Pot / i
				Man[index] = v
				if index == shooter {
					y++
				}
			}


	}
	if y == 0 {

		name = Clients[0]
		lastName = Clients[len(Clients) - 1]

		Clients[0] = lastName
        Clients[len(Clients) - 1] = name
        shooter = Clients[0]
	}
	Game.GameData.Pot = 0
	d.GameData.Pot = 0

	return M, Game, d, Clients, shooter
}

func matchedBet(M map[string]Bet) bool {
	i := 0
	t := 0


	for index, v := range M {


		if(M[index] != M[shooter]) {
		if v.Wager == M[shooter].Wager {
			t++
		}
	}
	}
	i++
	if t == 1 {
		return true
	}
	return false
}

func countBets(M map[string]Bet) bool {
	i := 0
	t := 0
	var id string

	for index, v := range M {
		if i == 0 {
			id = index
		}
		if v.Wager == M[id].Wager {
			t++

		}
		i++
	}

	if len(M) > 1 && t == len(M) {
		return true
	} else {
		return false
	}
}

func placeBet(M map[string]Bet, U []string, d Round) (map[string]Bet, []string, Round) {

	for index, v := range M {

		if index == d.UUID {

			if d.Wager != M[shooter].Wager && index != shooter {
				return M, U, d
			}
			v.Wager = d.Wager
			v.Bet = d.Message
			Game.GameData.Pot = Game.GameData.Pot + d.Wager
			d.GameData.Pot = Game.GameData.Pot
		    v.Balance = v.Balance - d.Wager
			M[d.UUID] = v
			U = append(U, d.UUID)
		}

	}
	
	return M, U, d
}

func walletReader(conn *websocket.Conn) {
for {
	_, p, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(p), &wallet); err != nil {
		panic(err)
	}

	fmt.Println(wallet)
}


}

func walletEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//log.Println("Client Successfully Connected...")

	walletReader(ws)


}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/v1/wallet", walletEndpoint)
    fmt.Println("Go Websockets!")
}

//redirect all HTTP traffic to HTTPS server on port :443
func redirectTLS(w http.ResponseWriter, r *http.Request)  {

	http.Redirect(w,r,"https://freshmintrecords.com:5005"+r.RequestURI,http.StatusMovedPermanently)
}

func main() {

	Man = make(map[string]Bet)
	Clients = make(map[int]string)
	setupRoutes()

    err := http.ListenAndServeTLS(":5000","./freshmintrecords_com.crt","./freshmintrecords.key",nil)
    if err != nil {
    	log.Fatal("ListenAndServe:", err)
    }
	//log.Fatal(http.ListenAndServe(":5000", nil))
}
