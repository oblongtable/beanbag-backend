package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	mywebsoc "github.com/oblongtable/beanbag-backend/websocket"
)

const (
	PORT  = "9091"
	DELAY = 3 * time.Second
)

var (
	server  *gin.Engine
	wssvr   *mywebsoc.WebSocServer
	client  *websocket.Conn
	client2 *websocket.Conn
	err     error
	res     *http.Response

	testRoomID string
)

func StartDummyServer() {
	server = gin.Default()
	wssvr = mywebsoc.NewWebSockServer()
	server.GET("/ws", wssvr.ServeWs)
	go server.Run(":" + PORT)
	time.Sleep(DELAY)
}

func TestUserRegister(t *testing.T) {
	client, res, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:9091/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	time.Sleep(DELAY)
	fmt.Println(wssvr.Clients)
	if len(wssvr.Clients) != 1 {
		t.Errorf("Client list size != 1; Expected 1; Current %d", len(wssvr.Clients))
	}
}

func TestUserCreateRoom(t *testing.T) {
	subevt := mywebsoc.CreateRoomEvent{
		RoomName: "Room1XDDD",
		RoomSize: 8,
	}

	rawSubEvtData, err := json.Marshal(subevt)
	if err != nil {
		log.Fatal("json marshal:", err)
	}

	evt := mywebsoc.Event{
		Type:    mywebsoc.EventCreateRoom,
		Payload: rawSubEvtData,
	}

	rawEvtData, err := json.Marshal(evt)
	if err != nil {
		log.Fatal("json marshal:", err)
	}
	fmt.Println("Event Triggerd:", string(rawEvtData))
	client.WriteMessage(websocket.TextMessage, rawEvtData)
	time.Sleep(DELAY)

	fmt.Println(wssvr.Rooms)
	if len(wssvr.Rooms) != 1 {
		t.Errorf("Room list size != 1; Expected 1; Current %d", len(wssvr.Rooms))
	}

	var evtCbMsg mywebsoc.EventCallbackMessage

	if msg_type, data, err2 := client.ReadMessage(); err2 != nil {
		t.Errorf("Message read failed: %v", err2)
	} else {
		fmt.Printf("%v %v %s\n", msg_type, data, string(data))
		err = json.Unmarshal(data, &evtCbMsg)
		if err != nil {
			t.Errorf("Unmarshal failed: %v", err)
		}

		var roomInfo mywebsoc.RoomInfo
		err = json.Unmarshal(evtCbMsg.Info, &roomInfo)
		if err != nil {
			t.Errorf("Unmarshal failed: %v", err)
		}

		testRoomID = roomInfo.ID
		fmt.Printf("Test Room ID: %s\n", testRoomID)
	}
}

func TestSingleUserJoinRoom(t *testing.T) {
	// client2, res, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:9091/ws", nil)
	// if err != nil {
	// 	log.Fatal("dial:", err)
	// }
	// time.Sleep(DELAY)
	// fmt.Println(wssvr.Clients)
	// if len(wssvr.Clients) != 2 {
	// 	t.Errorf("Client list size != 2; Expected 2; Current %d", len(wssvr.Clients))
	// }

	// subevt := mywebsoc.JoinRoomEvent{
	// 	RoomID: "Room1XDDD",
	// }

	// rawSubEvtData, err := json.Marshal(subevt)
	// if err != nil {
	// 	log.Fatal("json marshal:", err)
	// }

	// evt := mywebsoc.Event{
	// 	Type:    mywebsoc.EventCreateRoom,
	// 	Payload: rawSubEvtData,
	// }

	// rawEvtData, err := json.Marshal(evt)
	// if err != nil {
	// 	log.Fatal("json marshal:", err)
	// }
	// fmt.Println("Event Triggerd:", string(rawEvtData))
	// client.WriteMessage(websocket.TextMessage, rawEvtData)
	// time.Sleep(DELAY)

	// fmt.Println(wssvr.Rooms)
	// if len(wssvr.Rooms) != 1 {
	// 	t.Errorf("Room list size != 1; Expected 1; Current %d", len(wssvr.Rooms))
	// }
}

func TestSingleUserLeaveRoom(t *testing.T) {
	// TODO
}

func TestMultipleUserJoinLeaveRoom(t *testing.T) {
	// TODO
}

func TestMain(m *testing.M) {
	StartDummyServer()
	code := m.Run()
	os.Exit(code)
}
