package main

import (
	"log"
	"github.com/lastfreeacc/sinumbot/store"
)

func main() {

	// inline test for bolt storage
	botStore := store.NewBoltStore()
	botStore.SaveMemo(123, "memo1")
	botStore.SaveMemo(123, "memo2")
	u, _ := botStore.GetUser(123)
	log.Printf("%v", u)
	// ...

}