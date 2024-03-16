package database

import (
	"context"
	"fmt"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

type FireDb struct {
	*db.Client
	mu          sync.Mutex
	mockData    interface{} 
}

var fire FireDb


// Connecting with Realtime Database
func (db *FireDb) Connect() error {

	opt := option.WithCredentialsFile("E:/blog/database/private.json") // Private key for Authentication and Certificate

	
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.DatabaseWithURL(context.Background(),"https://blogbase-2269a-default-rtdb.europe-west1.firebasedatabase.app") // database link
	if err != nil {
		return fmt.Errorf("error initializing Firebase database client: %v", err)
	}
	db.Client = client
	return nil
}

func (db *FireDb) SetMockData(data interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.mockData = data
}

func FirebaseDB() *FireDb {
	if err := fire.Connect(); err != nil {
		panic(err)
	}
	return &fire
}
