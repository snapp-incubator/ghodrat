package ion_sfu

import "github.com/pion/webrtc/v3"

type Candidate struct {
	Target    int                  `json:"target"`
	Candidate *webrtc.ICECandidate `json:"candidate"`
}

type ResponseCandidate struct {
	Target    int                      `json:"target"`
	Candidate *webrtc.ICECandidateInit `json:"candidate"`
}

// SendOffer object to send to the sfu over Websockets.
type SendOffer struct {
	SID   string                     `json:"sid"`
	Offer *webrtc.SessionDescription `json:"offer"`
}

// SendAnswer object to send to the sfu over Websockets.
type SendAnswer struct {
	SID    string                     `json:"sid"`
	Answer *webrtc.SessionDescription `json:"answer"`
}

// TrickleResponse received from the sfu server.
type TrickleResponse struct {
	Params ResponseCandidate `json:"params"`
	Method string            `json:"method"`
}

// Response received from the sfu over Websockets.
type Response struct {
	Params *webrtc.SessionDescription `json:"params"`
	Result *webrtc.SessionDescription `json:"result"`
	Method string                     `json:"method"`
	ID     uint64                     `json:"id"`
}
