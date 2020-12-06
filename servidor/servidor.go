package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

var procesos = []int{1, 2, 3, 4, 5}
var contador_del_proceso uint64 = 0

type Informacion_del_proceso struct {
	Mensaje  string
	Proceso  int
	Contador uint64
}

func server() {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var msg Informacion_del_proceso
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	if msg.Mensaje == "Agregar proceso" {
		if len(procesos) > 0 {
			response := Informacion_del_proceso{
				Mensaje:     "Success",
				Proceso:  procesos[0],
				Contador: contador_del_proceso,
			}
			copy(procesos[0:], procesos[1:])       
			procesos = procesos[:len(procesos)-1] 
			err := gob.NewEncoder(c).Encode(response)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			response := Informacion_del_proceso{
				Mensaje:    "Error",
				Proceso: 0,
			}
			err := gob.NewEncoder(c).Encode(response)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		procesos = append(procesos, msg.Proceso)
	}
}

func printProccess() {
	go func() {
		for {
			contador_del_proceso++
			for i := 0; i < len(procesos); i++ {
				fmt.Println(procesos[i], "|", contador_del_proceso)
			}
			time.Sleep(time.Millisecond * 500)
			fmt.Println("----------")
		}
	}()
}

func main() {
	go server()
	printProccess()
	var input string
	fmt.Scanln(&input)
}
