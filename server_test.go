package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandleGetUser(t *testing.T) {
	server := NewServer(100);

	testServer := httptest.NewServer(http.HandlerFunc(server.handleGetUser));

	numberOfRequests := 1000;

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(numberOfRequests);

	for i := 0; i < numberOfRequests; i++ {
		go func(i int){
			id := i % 100 + 1;
			fmt.Println("id is ",id)

			url := fmt.Sprintf("%s/?id=%d", testServer.URL, id);

			response, err := http.Get(url);
			
			if err != nil {
				t.Error(err);
			}
			defer response.Body.Close();

			user := User{}

			if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
				t.Error(err);
			}
			fmt.Printf("%+v \n", user);
			waitGroup.Done();
		}(i);
		
		//a little sleep, because it's going as fast as our own laptop
		//in real condition, nothing makes http requests this much fast
		time.Sleep(time.Millisecond * 1);
	}

	waitGroup.Wait();

	fmt.Println("times we hit the database", server.dbhit);
}