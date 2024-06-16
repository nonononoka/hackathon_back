package main

import "back/controller"

func main() {
	controller.StartServer()
	//http.HandleFunc("/ws", websocket.NewWebsocketHandler().Handle)
	//
	//port := "80"
	//log.Printf("Listening on port %s", port)
	//if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
	//	log.Panicln("Serve Error:", err)
	//}
}
