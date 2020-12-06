package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

type Informacion_del_proceso struct {
	Mensaje  string
	Proceso  int
	Contador uint64
}

var Contador_del_proceso uint64 = 0
var Proceso_actual int = 0

func cliente(msg Informacion_del_proceso) {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	go handler_respuesta_servidor(c)
}

func handler_respuesta_servidor(c net.Conn) {
	var msg Informacion_del_proceso
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	} else {
		if msg.Mensaje == "Success" {
			Proceso_actual = msg.Proceso
			ejecutar_proceso(msg)
		} else {
			msg := Informacion_del_proceso{
				Mensaje:    "Agregar proceso",
				Proceso: 0,
			}
			go cliente(msg)
		}
	}
}

func devolver_proceso() {
	msg := Informacion_del_proceso{
		Mensaje:     "Devolver proceso",
		Proceso:  Proceso_actual,
		Contador: Contador_del_proceso,
	}
	cliente(msg)
}

func ejecutar_proceso(msg Informacion_del_proceso) {
	Contador_del_proceso = msg.Contador
	go func() {
		for {
			Contador_del_proceso++
			time.Sleep(time.Millisecond * 500)
			fmt.Println(msg.Proceso, "|", Contador_del_proceso)
		}
	}()
}

func main() {
	msg := Informacion_del_proceso{
		Mensaje: "Agregar proceso",
		Proceso: 0,
	}
	go cliente(msg)
	var input string
	fmt.Scanln(&input)
	devolver_proceso()
}
