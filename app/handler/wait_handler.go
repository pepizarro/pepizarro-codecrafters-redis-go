package handler

import (
	"fmt"
	"net"
)

func (h *Handler) WaitHandler(conn net.Conn, buffer []byte) {

	_, err := conn.Write(h.parser.WriteInteger(len(h.replicas)))
	if err != nil {
		fmt.Println("Error writing to client: ", err)
		return
	}
}
