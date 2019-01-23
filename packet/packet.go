package packet

import (
	"bytes"
	"database/sql"
	"kairos-go/remote"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/websocket"
)

type PacketName uint16

const (
	CSLogin PacketName = 10001

	SCLoggedIn PacketName = 20001
)

// Packet use for represent packet
type Packet struct{}

// ReceiveMessage use for read packet from websocket
func ReceiveMessage(packetID uint16, reader *bytes.Reader, remoteClient remote.Remote, clients map[*websocket.Conn]remote.Remote) {
	switch PacketName(packetID) {
	case CSLogin:
		log.Println("CSLogin")

		packetReader := Reader{
			BytesReader: reader,
		}

		name := packetReader.ReadString()
		number := packetReader.ReadInt8()

		log.Println(name)
		log.Println(number)

		db, err := sql.Open("mysql", "root:digitopolis@tcp(128.199.230.100:3306)/kairos")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// var (
		// 	nameFromDB string
		// 	age        int
		// )

		rows, err := db.Query("SELECT SLEEP(10.0)")
		// rows, err := db.Query("SELECT name, age FROM user WHERE name = ?", "armariya")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		// defer rows.Close()

		// for rows.Next() {
		// 	err := rows.Scan()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	// log.Println(nameFromDB, age)
		// }

		// err = rows.Err()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		remoteClient.Send(sendReceiveLoggedIn(), clients)
	}
}

func sendReceiveLoggedIn() []byte {
	var data = []interface{}{
		uint16(SCLoggedIn),
		"armariya",
		int8(20),
	}

	return WritePacket(data)
}
