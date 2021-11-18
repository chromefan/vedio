package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
)

var addr = flag.String("addr", "127.0.0.1:8089", "http service address")

const RTMP_SERVER = "rtmp://localhost:1935/hls/live"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)
		ffmpeg(message)
		//err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
func ffmpeg(msg []byte) {
	params := []string{
		"-re",
		"-i",
		"pipe:",
		"-vcodec",
		"copy",
		"-acodec",
		"aac",
		"-f",
		"flv",
		RTMP_SERVER,
	}
	in := bytes.NewBuffer(nil)
	cmd := exec.Command("/Users/luohuanjun/work/ffmpeg/ffmpeg", params...)
	fmt.Println(cmd.String(), params)
	cmd.Stdout = nil
	cmd.Stderr = nil
	in.Write(msg)
	cmd.Stdin = in
	err := cmd.Run()
	fmt.Println(cmd.String(), err)
}
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
