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

// GetUser ...
func (bs *boltStrore) GetUser(userID int64) (*User, error) {
	var v []byte
	byteID := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteID, uint64(userID))
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v = b.Get(byteID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		// user not found
		return nil, fmt.Errorf("user with id %d does not exist", userID)
	}
	u := User{}
	err = json.Unmarshal(v, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// SaveMemo ...
func (bs *boltStrore) SaveMemo(userID int64, memo *Memo) error {

	u, err := bs.GetUser(userID)
	if err != nil {
		log.Printf("[Info] can not get user with id: %d, err: %s\n", userID, err)
		ms := make([]*Memo, 0, 10)
		u = &User{ID: userID, Memos: ms}
	}
	u.Memos = append(u.Memos, memo)
	return bs.saveUser(u)
}

func (bs *boltStrore) saveUser(u *User) error {
	v, err := json.Marshal(u)
	if err != nil {
		return err
	}
	byteID := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteID, uint64(u.ID))
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		b.Put(byteID, v)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

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