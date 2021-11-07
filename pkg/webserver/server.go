package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/client/get_entries", server.getEntries)
	log.Println("listening :8069")
	log.Println(http.ListenAndServe(":8069", nil))
}

func (s *Server) jsonIngest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
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

const (
	queryStart = "start"
	queryEnd   = "end"
)

func (s *Server) getEntries(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	queryValues := req.URL.Query()
	var (
		start, end time.Time
	)
	startQuery := queryValues.Get(queryStart)
	if startQuery != "" {
		start, _ = time.Parse(time.RFC3339, startQuery)
	}
	endQuery := queryValues.Get(queryEnd)
	if endQuery != "" {
		end, _ = time.Parse(time.RFC3339, endQuery)
	}

	entries, err := s.entryStorer.GetRange(req.Context(), start, end)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	err = enc.Encode(entries)

	if err != nil {
		log.Println(err)
		return
	}
}
