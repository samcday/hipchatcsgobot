package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	steam "github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol"
	"github.com/Philipp15b/go-steam/protocol/steamlang"
	"github.com/Philipp15b/go-steam/steamid"
)

type debugHandler struct{}

func (*debugHandler) HandlePacket(p *protocol.Packet) {
	fmt.Printf("Steam packet %v", p.EMsg)
}

func main() {
	c := steam.NewClient()
	c.Connect()

	for event := range c.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			fmt.Println("Connected. Sending auth")
			c.Auth.LogOn(&steam.LogOnDetails{
				Username: os.Getenv("STEAM_USER"),
				Password: os.Getenv("STEAM_PASS"),
			})
		case *steam.MachineAuthUpdateEvent:
			fmt.Printf("Got auth hash yay %d\n", len(e.Hash))
		case *steam.ClientCMListEvent:
			fmt.Println("Got CM list")
			d, err := json.Marshal(e.Addresses)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile("servers.json", d, 0666)
			if err != nil {
				panic(err)
			}
		case *steam.LoggedOnEvent:
			fmt.Println("Logged on!")
			c.Social.SetPersonaState(steamlang.EPersonaState_Online)
		case *steam.FriendsListEvent:
			fmt.Printf("Friends? %d\n", c.Social.Friends.Count())
			id, err := steamid.NewId("STEAM_0:1:193009663")
			if err != nil {
				panic(err)
			}
			f, err := c.Social.Friends.ById(id)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(os.Stdout).Encode(f)
			fmt.Printf("Hmmm: %s\n", f.Name)
			if f.Relationship == steamlang.EFriendRelationship_RequestRecipient {
				c.Social.AddFriend(id)
				fmt.Println("Added friend!")
			}
		case *steam.FriendStateEvent:
			fmt.Printf("Friends? %d\n", c.Social.Friends.Count())
		case error:
			fmt.Println("Ohnoes error %v\n", e)
		}
	}
}
