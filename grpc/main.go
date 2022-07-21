package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	addr := `:8080`
	s, e := NewServer(addr)
	if e != nil {
		log.Fatalln(e)
	}
	log.Println(`listen on`, addr)

	e = s.Serve()
	if e != nil {
		log.Fatalln(e)
	}
}
