package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/codemonkeysoftware/mouseion/pkg/entry"
)

type Server struct {
	entryStorer EntryStorer
}

func New(entryStorer EntryStorer) *Server {
	return &Server{
		entryStorer: entryStorer,
	}
}

func (server *Server) Start() {
	http.HandleFunc("/ingest.json", server.jsonIngest)
	log.Println("listening :8069")
	log.Println(http.ListenAndServe(":8069", nil))
}

func (s *Server) jsonIngest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
	}
	defer func() {
		if err := req.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	entry := &entry.LogEntry{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(entry)
	if err != nil {
		log.Println(err)
	}

	err = s.entryStorer.Save(req.Context(), entry)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
	}
}
