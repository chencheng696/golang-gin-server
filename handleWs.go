package main

import (
	"fmt"
	"net/http"

	"Yinghao/klib"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/websocket"
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	// 必须死循环，gin通过协程调用该handler函数，一旦退出函数，ws会被主动销毁
	for {
		// recieve
		//mt, message, e := conn.ReadMessage()
		mt, _, e := conn.ReadMessage()
		if e != nil {
			fmt.Println("read: %+v", e)
			break
		}

		//fmt.Println("recv: %s", message)

		ret := make(map[string]interface{})
		ret[`ret`] = 0
		//ret["state_mq9"] = sensor_read_mq9()
		//ret["state_hcsr501"] = sensor_read_hcsr501()
		buf := []byte(klib.MapToJson(ret))

		e = conn.WriteMessage(mt, buf)
		if e != nil {
			fmt.Println("write: %+v", e)
			break
		}
	}
}
