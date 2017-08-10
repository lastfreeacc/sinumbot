package store

import (
	"fmt"
	"encoding/binary"
	"encoding/json"
	"log"
	"github.com/boltdb/bolt"
)

type boltStrore struct {
	db *bolt.DB
}


func (bs *boltStrore) GetUser(userID int64) (*User, error) {
	u := User{}
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		byteID := make([]byte, 8)
		binary.LittleEndian.PutUint64(byteID, uint64(userID))
		v := b.Get([]byte(byteID))
		if v == nil {
			log.Printf("[Info] user with id %d does not exist\n", userID)
			return fmt.Errorf("user with id %d does not exist", userID)
		}
		
		return json.Unmarshal(v, u)
	})
	return &u, err
}

func (bs *boltStrore) SaveMemo(userID int64, memo string) error {
	u, err := bs.GetUser(userID)
	if err != nil {
		log.Printf("[Info] can not get user, err: %s\n", err)
		ms := make([]string, 0, 10)
		u = &User{id: userID, memos: ms}
	}
	u.memos = append(u.memos, memo)
	
	// err := bs.db.Update(func(tx *bolt.Tx) error {
		
	// })
	return nil
}

//func (bs *boltStrore) createUser() 

// NewBoltStore ...
func NewBoltStore() Store {
	db, err := bolt.Open("sinumbot.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("users"))
        if err != nil {
            return err
		}
		return nil
    })
    if err != nil {
		log.Fatal(err)
	}
	return &boltStrore{db:db}
}