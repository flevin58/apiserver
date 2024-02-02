package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/flevin58/apiserver/router"
	"github.com/google/uuid"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	address       string
	router        *router.Router
	shoppingItems []Item
}

func NewServer() *Server {
	server := &Server{
		router:        router.New(),
		address:       ":8080",
		shoppingItems: []Item{},
	}
	server.setRoutes()
	return server
}

func (s *Server) WithAddress(address string) *Server {
	s.address = address
	return s
}

func (s *Server) Run() {
	log.Println("server listening for connections on address ", s.address)
	if err := http.ListenAndServe(s.address, s.router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed")
		}
		log.Fatalf("server failure: %v", err)
	}
}

func (s *Server) setRoutes() {

	s.router.AddRoute("/shopping-items", "GET", s.listShoppingItems)
	s.router.AddRoute("/shopping-items", "POST", s.createShoppingItem)
	s.router.AddRoute("/shopping-items/", "DELETE", s.removeShoppingItem)
}

func (s *Server) createShoppingItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = uuid.New()
	s.shoppingItems = append(s.shoppingItems, item)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) listShoppingItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) removeShoppingItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("deleting", id.String())
	for i, item := range s.shoppingItems {
		if item.ID == id {
			s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
			return
		}
	}

	// Here id not found
	http.Error(w, "item not found: "+idStr, http.StatusBadRequest)
}
