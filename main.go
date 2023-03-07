//lint:file-ignore ST1006 because self is a valid name

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	ID       int
	USername string
}

type Server struct {
	db map[int]*User
	cache map[int]*User
	dbhit int
}

func NewServer(count int) *Server {
	db := make(map[int]*User)

	for i := 0; i < count; i++ {
		db[i+1] = &User{
			ID:       i + 1,
			USername: fmt.Sprintf("user_%d", (i + 1)),
		}
	}

	return &Server {
		db: db,
		cache: make(map[int]*User),
	}
}

func (self *Server) tryCache(id int) (*User, bool){
	user, ok := self.cache[id];
	return user, ok;
}

func (self *Server) handleGetUser(responseWriter http.ResponseWriter, request *http.Request){
	idStr := request.URL.Query().Get("id");
	id, err := strconv.Atoi(idStr);
	if err != nil {
		log.Fatal("invalid request 😥");
		return;
	}

	//first, we try to hit the cache
	user, ok := self.cache[id];

	if ok {
		json.NewEncoder(responseWriter).Encode(user);
		return;
	}

	//hit the database
	user, ok = self.db[id];

	if !ok {
		panic("😱😱😱 user not found 😱😱😱")
	}

	self.dbhit++

	//insert in cache
	self.cache[id] = user;

	json.NewEncoder(responseWriter).Encode(user);
}

func main() {

}

