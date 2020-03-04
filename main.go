package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type MSG struct {
	Hash string
}

type Count struct {
	Size           int
	Percent        float64
	CompletedBlock int
	TotalBlock     int
	Speed          int
	Remain         int
}

func countServer(ws *websocket.Conn) {
	fmt.Println("LocalAddr", ws.LocalAddr())
	fmt.Println("RemoteAddr", ws.RemoteAddr())

	defer ws.Close()
	var msg MSG
	var count Count
	var before int
	var now int
	lk := sync.RWMutex{}
	frameTime := 3
	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		return
	}
	fmt.Println("开始处理:", msg.Hash)
	size := 104857600
	total := 612

	count.Size = size
	count.Percent = 0.00
	count.CompletedBlock = 0
	count.TotalBlock = total
	count.Speed = 0
	count.Remain = 0

	var cc int
	rand.Seed(time.Now().Unix())
	// TODO go feed
	for i := 0; i < 612; i++ {
		go func() {
			n := int64(rand.Intn(400))
			time.Sleep(time.Duration(n) * time.Second)
			lk.Lock()
			cc++
			lk.Unlock()
		}()

	}

	for {
		//TODO
		now = cc
		count.CompletedBlock = now
		count.Percent = float64(now) / float64(total)
		count.Speed = (now - before) * 256 * 1024 / frameTime
		if count.Speed == 0 {
			count.Remain = 0
		} else {
			count.Remain = (total - now) * 256 * 1024 / count.Speed
		}

		before = now
		fmt.Println(count)
		err = websocket.JSON.Send(ws, count)
		if err != nil {
			return
		}
		if count.CompletedBlock == count.TotalBlock {
			return
		}

		time.Sleep(time.Duration(frameTime) * time.Second)
	}
}

//
//func echoHandler(ws *websocket.Conn) {
//	var hash string
//	err := websocket.Message.Receive(ws, &hash) // &
//	if err != nil {
//
//	}
//	fmt.Println(hash)
//	_ = websocket.Message.Send(ws, strconv.Itoa(520))
//	fmt.Println("total:", strconv.Itoa(520))
//	i := 0
//	for {
//		i++
//		go func(num int) {
//			time.Sleep(1 * time.Second)
//			_ = websocket.Message.Send(ws, strconv.Itoa(num))
//			fmt.Println("ok:", strconv.Itoa(num))
//		}(i)
//
//		//time.Sleep(2 * time.Second)
//		if i == 10 {
//			time.Sleep(5 * time.Second)
//			_ = websocket.Message.Send(ws, "down")
//			break
//		}
//	}
//
//}

func main() {
	http.Handle("/count", websocket.Handler(countServer))
	//http.Handle("/echo", websocket.Handler(echoHandler))
	//http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
