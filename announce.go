package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Announce(changes chan Change, port uint) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		for {
			select {
			case change := <-changes:
				res.Header().Set("Access-Control-Allow-Origin", "*")
				res.Header().Set("Cache-Control", "nocache")
				res.Header().Set("Content-Type", "text/event-stream")
				if f, ok := res.(http.Flusher); ok {
					f.Flush()
				}

				b, _ := json.Marshal(change)
				msg := fmt.Sprintf("data: %s\n\n", b)
				res.Write([]byte(msg))
				if f, ok := res.(http.Flusher); ok {
					f.Flush()
				}

			case <-req.Context().Done():
				return
			}
		}
	})

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
