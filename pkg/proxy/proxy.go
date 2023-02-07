package proxy

import (
	"encoding/hex"
	"fmt"
	"github.com/instantminecraft/client/pkg"
	"github.com/instantminecraft/client/pkg/mcserver"
	"github.com/instantminecraft/client/pkg/server"
	"io"
	"log"
	"net"
)

const (
	PORT = 25585
)

func Start() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalln("Couldn't start tcp proxy at port", PORT, "=>", err)
	}
	defer l.Close()

	log.Println("Started proxy on port", PORT)
	log.Println("HTTP and Minecraft Clients should connect to this proxy")

	for {
		if conn, err := l.Accept(); err == nil {
			if pkg.DEBUG {
				log.Println("Accepting connection from", conn.RemoteAddr().String(), "...")
			}

			go acceptClient(conn)
		}
	}
}

func acceptClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2)
	_, err := conn.Read(buf)
	if err != nil {
		return
	}

	signature := hex.EncodeToString(buf)

	// select target port
	var targetPort int
	log.Println("Signature:", signature, "buffer:", buf)
	if isMinecraftConnection(signature) {
		// Proxy connection to minecraft server like nothing happened
		targetPort = mcserver.PORT
		if pkg.DEBUG {
			log.Println("Connection appears to be from a minecraft client. Redirecting...")
		}
	} else {
		// Proxy connection to local HTTP server
		targetPort = server.PORT
		if pkg.DEBUG {
			log.Println("Connection appears to be from a http client. Answering...")
		}
	}

	targetConnection, err := net.Dial("tcp", fmt.Sprintf(":%d", targetPort))
	if err != nil {
		// oh oh!
		return
	}

	// write the first read bytes
	if _, err = targetConnection.Write(buf); err != nil {
		return
	}

	// Start proxy
	go func() {
		if _, err := io.Copy(conn, targetConnection); err != nil {
			return
		}
	}()
	defer targetConnection.Close()
	if _, err := io.Copy(targetConnection, conn); err != nil {
		return
	}
}

func isMinecraftConnection(signature string) bool {
	// The first 4 digits area always "1000" or "1500" if a minecraft client tries to connect
	// Note: I don't know why, but on local running minecraft servers "1000" is the signature while on deployed servers the signatures seems to be "1500"
	return signature == "1000" || signature == "1500"
}
