package ion_sfu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

type Ion_sfu struct {
	Logger *zap.Logger
	Client *client.Client
	Config *Config

	connection   *websocket.Conn
	connectionID uint64
}

func (ion_sfu *Ion_sfu) dial() {
	addr := ion_sfu.Config.Address
	var err error

	ion_sfu.connection, _, err = websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
}

func (ion_sfu *Ion_sfu) onIceCandidate(candidate *webrtc.ICECandidate) {
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

	message := &jsonrpc2.Request{
		Method: "trickle",
		Params: params,
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(message)

	messageBytes := reqBodyBytes.Bytes()
	ion_sfu.connection.WriteMessage(websocket.TextMessage, messageBytes)
}

func (ion_sfu *Ion_sfu) offer() {
	offerJSON, err := json.Marshal(&SendOffer{
		Offer: ion_sfu.Client.GetLocalDescription(),
		SID:   "test room",
	})

	if err != nil {
		panic(err)
	}

	params := (*json.RawMessage)(&offerJSON)

	ion_sfu.connectionID = uint64(uuid.New().ID())

	offerMessage := &jsonrpc2.Request{
		Method: "join",
		Params: params,
		ID: jsonrpc2.ID{
			Num:      ion_sfu.connectionID,
			IsString: false,
			Str:      "",
		},
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(offerMessage)

	// send the offer over to the sfu using Websockets
	messageBytes := reqBodyBytes.Bytes()
	ion_sfu.connection.WriteMessage(websocket.TextMessage, messageBytes)
}

func (ion_sfu *Ion_sfu) readMessage(connection *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		_, message, err := connection.ReadMessage()
		if err != nil || err == io.EOF {
			log.Fatal("Error reading: ", err)
			break
		}

		fmt.Printf("recv: %s", message)

		var response Response
		json.Unmarshal(message, &response)

		// determine which event the message is for and handle them accordingly
		if response.Id == ion_sfu.connectionID {
			result := *response.Result
			if err := peerConnection.SetRemoteDescription(result); err != nil {
				log.Fatal(err)
			}
		} else if response.Id != 0 && response.Method == "offer" {
			// the sfu sends an offer and we react by saving the send offer into the remote
			// description of our peer connection and sending back an answer with the
			// local description so we can connect to the remote peer.

			peerConnection.SetRemoteDescription(*response.Params)
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				log.Fatal(err)
			}

			peerConnection.SetLocalDescription(answer)

			connectionUUID := uuid.New()
			ion_sfu.connectionID = uint64(connectionUUID.ID())

			offerJSON, err := json.Marshal(&SendAnswer{
				Answer: peerConnection.LocalDescription(),
				SID:    "test room",
			})
			if err != nil {
				log.Fatal(err)
			}

			params := (*json.RawMessage)(&offerJSON)

			answerMessage := &jsonrpc2.Request{
				Method: "answer",
				Params: params,
				ID: jsonrpc2.ID{
					IsString: false,
					Str:      "",
					Num:      ion_sfu.connectionID,
				},
			}

			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(answerMessage)

			messageBytes := reqBodyBytes.Bytes()
			connection.WriteMessage(websocket.TextMessage, messageBytes)
		} else if response.Method == "trickle" {
			// The sfu sends a new ICE candidate and we add it to the peer connection

			var trickleResponse TrickleResponse

			if err := json.Unmarshal(message, &trickleResponse); err != nil {
				log.Fatal(err)
			}

			err := peerConnection.AddICECandidate(*trickleResponse.Params.Candidate)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
