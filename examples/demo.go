package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bloomapi/go-dailyco"
)

func main() {
	c := dailyco.Client{AuthorizationToken: os.Getenv("DAILYCO_TOKEN")}

	room, err := c.CreateRoom("", "private", &dailyco.RoomProperties{})
	if room != nil {
		fmt.Println(room.Id)
		fmt.Println(room.Name)
		fmt.Println(room.CreatedAt)
	} else {
		fmt.Println(err)
		return
	}

	expiration := time.Now().UTC().AddDate(0, 0, 1).Unix()
	roomName := room.Name
	response, err := c.CreateMeetingToken(&dailyco.MeetingTokenProperties{
		RoomName: &roomName,
		Exp:      &expiration,
	})

	if response != nil {
		fmt.Println(response.Token)
	}

	result, err := c.DeleteRoom(room.Name)
	fmt.Printf("%t", result)
	fmt.Println(err)
}
