package storescp

import (
	"fmt"
	"net"
)

func Echocheat() {
	listener, err := net.Listen("tcp", ":11112")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection2(conn)
	}
}

func handleConnection2(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	buf2 := make([]byte, 512)
	buf3 := make([]byte, 512)
	buf4 := make([]byte, 512)
	buf5 := make([]byte, 512)
	buf6 := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		for {
			conn2, err := net.Dial("tcp", "localhost:11113")
			if err != nil {
				panic(err)
			}
			conn2.Write(buf[:n])
			n2, err := conn2.Read(buf2)
			if err != nil {
				panic(err)
			}
			// Below is not needed, just for debugging
			// fmt.Println("buf2", buf2[:n2])
			//trim the message
			buf2 = trimMessage(buf2)
			//decode the message
			AARQStruct, err := DecodeAAssociateRQ(buf2)
			if err != nil {
				panic(err)
			}
			fmt.Println(AARQStruct.ToString())

			conn.Write(buf2[:n2])
			n3, err := conn.Read(buf3)
			if err != nil {
				panic(err)
			}
			conn2.Write(buf3[:n3])
			n4, err := conn2.Read(buf4)
			if err != nil {
				panic(err)
			}
			conn.Write(buf4[:n4])
			n5, err := conn.Read(buf5)
			if err != nil {
				panic(err)
			}
			conn2.Write(buf5[:n5])

			n6, err := conn2.Read(buf6)
			if err != nil {
				panic(err)
			}
			conn.Write(buf6[:n6])

			println(string(buf[:n]))
			println(string(buf2[:n2]))
			println(string(buf3[:n3]))
			println(string(buf4[:n4]))
			println(string(buf5[:n5]))
			println(string(buf6[:n6]))

		}
	}
}
