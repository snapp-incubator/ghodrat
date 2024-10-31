package ion_sfu

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

// nolint: gochecknoglobals
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type IonSfu struct {
	Logger *zap.Logger
	Client *client.Client
	Config *Config

	connection   *websocket.Conn
	connectionID uint64
	sid          string
}

func (ionSfu *IonSfu) dial() {
	addr := ionSfu.Config.Address

	var err error

	ionSfu.connection, _, err = websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
}

func (ionSfu *IonSfu) generateSID() {
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	ionSfu.sid = string(b)
}

func (ionSfu *IonSfu) onIceCandidate(candidate *webrtc.ICECandidate) {
	if candidate == nil {
		return
	}

	candidateJSON, err := json.Marshal(&Candidate{
		Candidate: candidate,
		Target:    0,
	})

	params := (*json.RawMessage)(&candidateJSON)

	if err != nil {
		log.Fatal(err)
	}

	// nolint: exhaustruct
	message := &jsonrpc2.Request{
		Method: "trickle",
		Params: params,
	}

	reqBodyBytes := new(bytes.Buffer)

	err = json.NewEncoder(reqBodyBytes).Encode(message)
	if err != nil {
		log.Fatal("cannot encode message: %", err)
	}

	messageBytes := reqBodyBytes.Bytes()

	err = ionSfu.connection.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Fatal("cannot write message: %", err)
	}
}

func (ionSfu *IonSfu) offer() {
	offerJSON, err := json.Marshal(&SendOffer{
		Offer: ionSfu.Client.GetLocalDescription(),
		SID:   ionSfu.sid,
	})
	if err != nil {
		panic(err)
	}

	params := (*json.RawMessage)(&offerJSON)

	ionSfu.connectionID = uint64(uuid.New().ID())

	// nolint: exhaustruct
	offerMessage := &jsonrpc2.Request{
		Method: "join",
		Params: params,
		ID: jsonrpc2.ID{
			Num:      ionSfu.connectionID,
			IsString: false,
			Str:      "",
		},
	}

	reqBodyBytes := new(bytes.Buffer)

	err = json.NewEncoder(reqBodyBytes).Encode(offerMessage)
	if err != nil {
		log.Fatal("cannot encode message: %", err)
	}

	// send the offer over to the sfu using Websockets
	messageBytes := reqBodyBytes.Bytes()

	err = ionSfu.connection.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Fatal("cannot write message: %", err)
	}
}

func (ionSfu *IonSfu) readMessage() {
	for {
		_, message, err := ionSfu.connection.ReadMessage()
		if err != nil || errors.Is(err, io.EOF) {
			log.Fatal("Error reading: ", err)
		}

		log.Printf("\nrecv: %s\n", message)

		var response Response

		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Fatal("error marshaling json: ", err)
		}

		// determine which event the message is for and handle them accordingly
		// nolint: nestif
		if response.ID == ionSfu.connectionID {
			ionSfu.Client.SetRemoteDescription(*response.Result)
		} else if response.ID != 0 && response.Method == "offer" {
			// the sfu sends an offer and we react by saving the send offer into the remote
			// description of our peer connection and sending back an answer with the
			// local description so we can connect to the remote peer.
			ionSfu.Client.SetRemoteDescription(*response.Result)
			ionSfu.Client.CreateAndSetAnswer()

			connectionUUID := uuid.New()
			ionSfu.connectionID = uint64(connectionUUID.ID())

			offerJSON, err := json.Marshal(&SendAnswer{
				Answer: ionSfu.Client.GetLocalDescription(),
				SID:    ionSfu.sid,
			})
			if err != nil {
				log.Fatal(err)
			}

			params := (*json.RawMessage)(&offerJSON)

			// nolint: exhaustruct
			answerMessage := &jsonrpc2.Request{
				Method: "answer",
				Params: params,
				ID: jsonrpc2.ID{
					Num:      ionSfu.connectionID,
					IsString: false,
					Str:      "",
				},
			}

			reqBodyBytes := new(bytes.Buffer)

			err = json.NewEncoder(reqBodyBytes).Encode(answerMessage)
			if err != nil {
				log.Fatal("cannot encode json: ", err)
			}

			messageBytes := reqBodyBytes.Bytes()

			err = ionSfu.connection.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				log.Fatal("cannot write message to socket: %w", err)
			}
		} else if response.Method == "trickle" {
			// The sfu sends a new ICE candidate and we add it to the peer connection
			var trickleResponse TrickleResponse

			if err := json.Unmarshal(message, &trickleResponse); err != nil {
				log.Fatal(err)
			}

			ionSfu.Client.AddIceCandidate(trickleResponse.Params.Candidate)
		}
	}
}
