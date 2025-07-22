package serv

import (
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	mux  *http.ServeMux
	port string
}

func (s *Service) Run() {
	log.Println("Start http server on :" + s.port)
	http.HandleFunc("/incoming", func(w http.ResponseWriter, req *http.Request) {

	})

	http.HandleFunc("/alive", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	})

	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		panic(err)
	}
}

func NewService(port string) (*Service, error) {
	s := &Service{
		port: port,
	}
	return s, nil
}
