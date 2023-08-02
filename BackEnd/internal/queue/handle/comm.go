package handle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// JSON containing the data,
// for more information see
// the documentation.
type NestedJSON struct {
	Id        int         `json:"id"`
	MsgType   int         `json:"msgtype"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// JSON where the various NestedJSON are packed.
type MessageJSON struct {
	Id    int          `json:"id"`
	Ndata int          `json:"ndata"`
	Data  []NestedJSON `json:"datajson"`
}

// This function decodes the payload.
// Returns the number of data received, the NestedJSON structure and the error.
func decodePayload(payload []byte) (int, []NestedJSON, error) {
	var id int
	var msgType int
	var value interface{}
	var dataJSON []NestedJSON
	nData := 0
	buf := bytes.NewReader(payload)
	var err error
	for {

		// Gets the transmission ID in byte format.
		id_bytes := make([]byte, 4)
		if err := binary.Read(buf, binary.BigEndian, &id_bytes); err != nil {
			return nData, dataJSON, err
		}
		if id, err = strconv.Atoi(string(id_bytes[2:])); err != nil {
			return nData, dataJSON, err
		}

		// Gets the type of data in byte format.
		// 0: bool
		// 1: Int
		// 2: Float
		// 3: String
		msgType_bytes := make([]byte, 6)
		if err := binary.Read(buf, binary.BigEndian, &msgType_bytes); err != nil {
			return nData, dataJSON, err
		}
		if msgType, err = strconv.Atoi(string(msgType_bytes[4:])); err != nil {
			return nData, dataJSON, err
		}

		// Stores the data in an interface.
		switch msgType {
		case 0: // Bool
			v := make([]byte, 5)
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nData, dataJSON, err
			}
			value = string(v[4:])

		case 1: // Int
			v := make([]byte, 8)
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nData, dataJSON, err
			}
			value = string(v[4:])

		case 2: // Float
			v := make([]byte, 9)
			if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
				return nData, dataJSON, err
			}
			value = string(v[4:])

		case 3: // String
			lenght_byte := make([]byte, 4)
			if err := binary.Read(buf, binary.BigEndian, &lenght_byte); err != nil {
				return nData, dataJSON, err
			}

			l, err := strconv.Atoi(string(lenght_byte[:]))
			if err != nil {
				return nData, dataJSON, err
			}
			stringBuf := make([]byte, l)
			if _, err := io.ReadFull(buf, stringBuf); err != nil {
				return nData, dataJSON, err
			}
			value = string(stringBuf)

		default:
			return nData, dataJSON, fmt.Errorf("unknown data type: %d", msgType)
		}

		// Creates the JSON structure with the obtained data.
		data := NestedJSON{Id: id, MsgType: msgType, Data: value, Timestamp: time.Now().Unix()}
		dataJSON = append(dataJSON, data)
		nData++
		if buf.Len() < 1 {
			break
		}
	}

	return nData, dataJSON, nil
}

// This function listens to the corresponding port and obtains the payload data in the final format.
func GetValue(messages chan<- MessageJSON) {

	addrip := os.Getenv("ADDRIP")
	port := os.Getenv("PORT")
	log.Println(" > addrip = " + addrip + ":" + port)

	// Listening over TCP connections
	listener, err := net.Listen("tcp", addrip+":"+port)
	if err != nil {
		log.Fatalf("Error when listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf(" > Error accepting connection: %v", err)
			continue
		}

		// Process the connection in a goroutine to keep accepting other connections
		go handleConnection(conn, messages)
	}
}

// This function is in charge of reading the buffer obtained from the hardware and
// returning in a channel the JSON ready to be enqueued.
func handleConnection(conn net.Conn, messages chan<- MessageJSON) {
	defer conn.Close()
	idPayload := 0
	for {
		// Define a buffer to store the read bytes
		length := make([]byte, 4) // Size 4, adjust the size according to your needs

		// Read data type 'type' from payload
		if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
			if err != io.EOF {
				log.Printf(" > Error lenght when reading data type: %v", err)
			}
			break
		}
		// Get length to make byte buffer
		l, _ := strconv.Atoi(string(length[:]))

		// Define a buffer to store the read bytes
		buf := make([]byte, l) // Size 4, adjust the size according to your needs

		// Read the data type 'type' from the payload
		if err := binary.Read(conn, binary.BigEndian, &buf); err != nil {
			if err != io.EOF {
				log.Printf(" > Error reading data type: %v", err)
			}
			break
		}

		// Get data from tcp payload
		nData, data, err := decodePayload(buf)
		if err != nil {
			log.Printf(" > Error decodePayload: %v", err)
			break
		}

		// Here you process the received message
		MessageJSON := MessageJSON{
			Id:    idPayload,
			Ndata: nData,
			Data:  data,
		}
		idPayload++
		messages <- MessageJSON
	}
}
