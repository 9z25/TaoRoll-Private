package main

import (
	"encoding/json"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
	"honnef.co/go/js/dom"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

type Bet int32

const (
	PASS     Bet = 0
	DONTPASS Bet = 1
)

func addOneDollar(balance int, wager int, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement, wa *dom.HTMLHeadingElement, ba *dom.HTMLHeadingElement) (int, int) {

	if balance > 0 {
		balance = balance - 1
		wager++
		wa.Set("innerText", wager)
		
	}

	if wager > 0 {
		p.Disabled = false
		dp.Disabled = false
	}
	return balance, wager
}
func addFiveDollars(balance int, wager int, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement, wa *dom.HTMLHeadingElement, ba *dom.HTMLHeadingElement) (int, int) {

	if balance > 4 {
		balance = balance - 5
		wager = wager + 5
		wa.Set("innerText", wager)
		
	}

	if wager > 0 {
		p.Disabled = false
		dp.Disabled = false
	}

	return balance, wager
}

func addTwentyFiveDollars(balance int, wager int, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement, wa *dom.HTMLHeadingElement, ba *dom.HTMLHeadingElement) (int, int) {
	if balance > 24 {
		balance = balance - 25
		wager = wager + 25
		wa.Set("innerText", wager)
		
	}

	if wager > 0 {
		p.Disabled = false
		dp.Disabled = false
	}
	return balance, wager
}

func addOneHundo(balance int, wager int, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement, wa *dom.HTMLHeadingElement, ba *dom.HTMLHeadingElement) (int, int) {
	if balance > 99 {
		balance = balance - 100
		wager = wager + 100
		wa.Set("innerText", wager)
		
	}

	if wager > 0 {
		p.Disabled = false
		dp.Disabled = false
	}
	return balance, wager
}

func playPass(play Bet, state GameState, wager int, b *dom.HTMLButtonElement, oneDollar *dom.HTMLButtonElement, fiveDollar *dom.HTMLButtonElement, twentyFiveDollar *dom.HTMLButtonElement, oneHundo *dom.HTMLButtonElement, result *dom.HTMLHeadingElement, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement) Bet {
	play = PASS



	if wager > 0 && state == COMEOUT {

		oneDollar.Disabled = true
		fiveDollar.Disabled = true
		twentyFiveDollar.Disabled = true
		oneHundo.Disabled = true

		p.Disabled = true
		dp.Disabled = true

	}

	return play
}

func playDontPass(play Bet, state GameState, wager int, b *dom.HTMLButtonElement, oneDollar *dom.HTMLButtonElement, fiveDollar *dom.HTMLButtonElement, twentyFiveDollar *dom.HTMLButtonElement, oneHundo *dom.HTMLButtonElement, result *dom.HTMLHeadingElement, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement) Bet {
	play = DONTPASS
	state = COMEOUT


	if wager > 0 && state == COMEOUT {

		oneDollar.Disabled = true
		fiveDollar.Disabled = true
		twentyFiveDollar.Disabled = true
		oneHundo.Disabled = true

		p.Disabled = true
		dp.Disabled = true
	}

	return play
}

func loadWager(result *dom.HTMLHeadingElement, w *dom.HTMLButtonElement, p *dom.HTMLButtonElement, dp *dom.HTMLButtonElement, oneDollar *dom.HTMLButtonElement, fiveDollar *dom.HTMLButtonElement, twentyFiveDollar *dom.HTMLButtonElement, oneHundo *dom.HTMLButtonElement, state GameState) GameState {

	state = COMEOUT
	w.Disabled = true

	if state == COMEOUT {
		oneDollar.Disabled = false
		fiveDollar.Disabled = false
		twentyFiveDollar.Disabled = false
		oneHundo.Disabled = false

		p.Disabled = true
		dp.Disabled = true
	}
	return state
}

func roll(i int, f interface{}) (int) {

	if i > 4 {
		js.Global.Call("clearInterval", f)
	}

	return i
}


	


type Player struct {
	Username string `json:"username,omitempty"`
	Bet      string `json:"bet,omitempty"`
}

type Players struct {
	Players  []Player `json:"players,omitempty"`
	Pot      int      `json:"pot,omitempty"`
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

type Dice struct {
	L int `json:"l,omitempty"`
	R int `json:"r,omitempty"`
}

type Data struct {
	*js.Object
	Users []string `js:"users,omitempty"`
}

type WalletJSON {
	Action string `json:"action"`
	Address string `json:"address"`
	Withdraw string `json:"withdraw"`
	Recipient string `json:"recipient"`
	Balance string `json:"balance"`
	UUID    string `json:"uuid"`
}

var w walletJSON
var g Round
var c bool
var Guid string

func main() {

	//for dom
	var f interface{}
	var leftDice int
	var rightDice int

	T := js.Global
	dice_arr := [6]string{"./res/dice1.png", "./res/dice2.png", "./res/dice3.png", "./res/dice4.png", "./res/dice5.png", "./res/dice6.png"}

	//for game
	var t int
	var point int
	var wager int
	var balance int
	wager = 0
	balance = 100

	d := dom.GetWindow().Document()

	//screens
	ls := d.GetElementByID("login-screen").(*dom.HTMLDivElement)
	gs := d.GetElementByID("gamble-screen").(*dom.HTMLDivElement)

	//dice
	l := d.GetElementByID("left-dice").(*dom.HTMLImageElement)
	r := d.GetElementByID("right-dice").(*dom.HTMLImageElement)

	//chips
	oneDollar := d.GetElementByID("one-dollar").(*dom.HTMLButtonElement)
	fiveDollar := d.GetElementByID("five-dollar").(*dom.HTMLButtonElement)
	twentyFiveDollar := d.GetElementByID("twentyfive-dollar").(*dom.HTMLButtonElement)
	oneHundo := d.GetElementByID("onehundred-dollar").(*dom.HTMLButtonElement)

	//buttons
	sg := d.GetElementByID("enter").(*dom.HTMLButtonElement)
	b := d.GetElementByID("roll-button").(*dom.HTMLButtonElement)
	w := d.GetElementByID("wager").(*dom.HTMLButtonElement)
	p := d.GetElementByID("pass").(*dom.HTMLButtonElement)
	dp := d.GetElementByID("dont-pass").(*dom.HTMLButtonElement)

	
	p.Disabled = true
	dp.Disabled = true
	w.Disabled = false

	//textboxes
	un := d.GetElementByID("username").(*dom.HTMLInputElement)
	result := d.GetElementByID("result").(*dom.HTMLHeadingElement)
	pb := d.GetElementByID("point-box").(*dom.HTMLHeadingElement)
	wa := d.GetElementByID("wager-amount").(*dom.HTMLHeadingElement)
	ba := d.GetElementByID("balance-amount").(*dom.HTMLHeadingElement)
	pa := d.GetElementByID("pot-amount").(*dom.HTMLHeadingElement)
	dr := d.GetElementByID("dice-roll").(*dom.HTMLHeadingElement)
	//rb := d.GetElementByID("room").(*dom.HTMLInputElement)

	//wallet textboxes 
	
	//sa := d.GetElementByID("sendAddress").(*dom.HTMLInputElement)
	da := d.GetElementByID("depositAddress").(*dom.HTMLInputElement)
	//wallet
	mw := d.GetElementByID("myWallet").(*dom.HTMLDivElement)
	
	//wallet buttons
	ow := d.GetElementByID("open-wallet").(*dom.HTMLButtonElement)
	cw := d.GetElementByID("close-wallet").(*dom.HTMLSpanElement)
	cp := d.GetElementByID("copy-address").(*dom.HTMLButtonElement)
	wm := d.GetElementByID("withdraw-money").(*dom.HTMLButtonElement)

	wa.Set("innerText", 0)

	var state GameState
	var play Bet
    var hostName string

	hostName = js.Global.Get("window").Get("location").Get("hostname").String()
	
	fmt.Print(hostName)
	gs.Class().Add("invisible")

	fmt.Println("Connection attempt ...")
	
	

	ws, err := websocketjs.New("wss://"+"freshmintrecords.com"+":5000/ws") // Does not block.
	if err != nil {
		// handle error
		fmt.Println(err)

	}

	//	m, err := websocketjs.New("ws://localhost:8080/v1/ws") // Does not block.

	//m, err := websocketjs.New("ws://localhost:8080/v1/ws")

	onOpen := func(ev *js.Object) {

		fmt.Println("Connection success!!")

	}
	// ...

	onMessage := func(ev *js.Object) {
		js.Global.Get("document").Call("querySelector", "#room").Set("innerHTML", "")

		Round := js.Global.Get("JSON").Call("parse", ev.Get("data"))
		fmt.Print(ev.Get("data"), "ln 367")

		diceL := Round.Get("gameData").Get("dice").Get("l").Int()
		diceR := Round.Get("gameData").Get("dice").Get("r").Int()
		num := Round.Get("gameData").Get("dice").Get("total").Int()
		shooter := Round.Get("gameData").Get("shooter").String()

		result.Set("innerText", Round.Get("jumbotron").String())


		l.Set("src", dice_arr[diceL])
		r.Set("src", dice_arr[diceR])

		

		dr.Set("innerText", num)


		usr := Round.Get("users")
		roller := Round.Get("gameData").Get("roll")
		//jsonState := Round.Get("state").String()
		wager := Round.Get("gameData").Get("wager").String()
		jsonPoint := Round.Get("gameData").Get("point").String()
		balanceAmount := Round.Get("gameData").Get("balance").String()
		potAmount := Round.Get("gameData").Get("pot").String()


		


		

		if wager != "undefined" {
		wa.Set("innerText",wager)
	}

		if potAmount == "undefined"{
			pa.Set("innerText", "$0")
		} else {
			pa.Set("innerText", "$" + potAmount)
		}
		

		if balanceAmount == "undefined" {
			ba.Set("innerText", "$0")
		} else {
			ba.Set("innerText","Balance: $" + balanceAmount)
		}
		if jsonPoint != "undefined" {
			pb.Set("innerText", jsonPoint)
		}
		
		if state == WAGER || roller.String() != "true" {
			w.Disabled = true
			b.Disabled = true
		}

		if Round.Get("gameData").Get("placeBet").String() == "true" {
			w.Disabled = false
		}

		if roller.String() == "true" && state == COMEOUT {
			b.Disabled = false
		}
		count := 0

		if usr.String() != "undefined" {
		var splits = strings.Split(usr.String(), ",")

		for _, element := range splits {
			if element != "" {
				li := js.Global.Get("document").Call("createElement", "li")
				span := js.Global.Get("document").Call("createElement", "span")
				if shooter == element {
					span.Get("classList").Call("add", "glyph-magic")
				}
				count = count + 1
				text := js.Global.Get("document").Call("createTextNode", element)
				li.Call("appendChild", text)
				li.Call("appendChild", span)
				js.Global.Get("document").Call("querySelector", "#room").Call("appendChild", li)
			}
		}
	}

	}

	onClose := func(ev *js.Object) {
		fmt.Sprintf("%b", ev)
		ws.Object.Call("close") // Send a text frame.
		// ...
		fmt.Println("Connection closed")
	}

	onError := func(ev *js.Object) {
		fmt.Println("Error", ev.Get("code").Int())
		// ...
	}
	//m.AddEventListener("open", false, loadRound)
	//m.AddEventListener("message", false, fetchUsers)
	ws.AddEventListener("open", false, onOpen)

	ws.AddEventListener("close", false, onClose)
	ws.AddEventListener("error", false, onError)

	ws.AddEventListener("message", false, onMessage)

	//login
	sg.AddEventListener("click", false, func(event dom.Event) {
		b := make([]byte, 16)

		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

		Guid = uuid

		if un.Value != "" {
			gs.Class().Toggle("invisible")
			ls.Class().Add("invisible")
fmt.Println("sg")
			a := &Round{Users: []string{}, GameData: Players{Players: []Player{Player{Username: "", Bet: ""}}}, UUID: Guid, Message: string(un.Value)}
			e, err := json.Marshal(a)
			if err != nil {
				fmt.Println(err)
				return
			}

			ws.Object.Call("send", e)

		}
	})

	// BEGIN

	rand.Seed(time.Now().UnixNano())

	state = WAGER

	w.AddEventListener("click", false, func(event dom.Event) {
		state = loadWager(result, w, p, dp, oneDollar, fiveDollar, twentyFiveDollar, oneHundo, state)
	})

	p.AddEventListener("click", false, func(event dom.Event) {

		amount, err := strconv.Atoi(wa.Get("innerText").String())
		if err != nil {
			fmt.Println("err")
		}
		a := &Round{Users: []string{}, GameData: Players{Players: []Player{Player{}}}, Wager: amount, UUID: Guid, Message: "PASS"}
		e, err := json.Marshal(a)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("p")
		ws.Object.Call("send", e)
		play = playPass(play, state, wager, b, oneDollar, fiveDollar, twentyFiveDollar, oneHundo, result, p, dp)
	})

	dp.AddEventListener("click", false, func(event dom.Event) {
		amount, err := strconv.Atoi(wa.Get("innerText").String())
		if err != nil {
			fmt.Println("err")
		}
		a := &Round{Users: []string{}, GameData: Players{Players: []Player{Player{}}}, Wager: amount, UUID: Guid, Message: "DONTPASS"}
		e, err := json.Marshal(a)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("dp")
		ws.Object.Call("send", e)
		play = playDontPass(play, state, wager, b, oneDollar, fiveDollar, twentyFiveDollar, oneHundo, result, p, dp)
	})

	i := 0

	oneDollar.Disabled = true
	fiveDollar.Disabled = true
	twentyFiveDollar.Disabled = true
	oneHundo.Disabled = true

	pb.Set("innerText", "")

	//wallet
	ow.AddEventListener("click", false, func(event dom.Event) {
		mw.Style().SetProperty("display","block","")
	})

	cw.AddEventListener("click", false, func(event dom.Event) {
		mw.Style().SetProperty("display","none","")
	})

	cp.AddEventListener("click", false, func(event dom.Event) {
		js.Global.Call("alert", da.Value)
	})

	wm.AddEventListener("click", false, func(event dom.Event) {
		js.Global.Call("alert", da.Value)
	})
    

	//gamble

	oneDollar.AddEventListener("click", false, func(event dom.Event) {
		balance, wager = addOneDollar(balance, wager, p, dp, wa, ba)
	})

	fiveDollar.AddEventListener("click", false, func(event dom.Event) {
		balance, wager = addFiveDollars(balance, wager, p, dp, wa, ba)
	})

	twentyFiveDollar.AddEventListener("click", false, func(event dom.Event) {
		balance, wager = addTwentyFiveDollars(balance, wager, p, dp, wa, ba)
	})

	oneHundo.AddEventListener("click", false, func(event dom.Event) {
		balance, wager = addOneHundo(balance, wager, p, dp, wa, ba)
	})



	b.AddEventListener("click", false, func(event dom.Event) {

		f = T.Call("setInterval", func() {
			leftDice = rand.Intn(6)
			rightDice = rand.Intn(6)
			b.Disabled = true
			//fmt.Println(i)
			t = leftDice + rightDice + 2
		fmt.Println("rolltotal",t)		
			

			i = roll(i, f)

	
if i == 5 {
			
			l.Set("src", dice_arr[leftDice])
			r.Set("src", dice_arr[rightDice])
fmt.Println("finishDice",leftDice,rightDice)
fmt.Println("5total",t)		
fmt.Println("5point",point)
					a := &Round{Users: []string{}, GameData: Players{Dice: Dice{L: leftDice, R: rightDice}}, UUID: Guid, Message: "finished"}
					e, err := json.Marshal(a)
					fmt.Println(a)
					if err != nil {
						fmt.Println(err)
						
					}

					ws.Object.Call("send", e)
					i = 0
				}
				
				
				if i < 5 && i != 0 {
					fmt.Println("rollingDice",leftDice,rightDice)
					a := &Round{Users: []string{}, GameData: Players{Dice: Dice{L: leftDice, R: rightDice}}, UUID: Guid, Message: "rolling"}
					
					e, err := json.Marshal(a)
					if err != nil {
						fmt.Println(err)
						
					}
					ws.Object.Call("send", e)
				}

				
			i++	

		}, 65)


	})
}
