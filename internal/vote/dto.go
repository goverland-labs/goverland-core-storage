package vote

import "encoding/json"

type ValidateRequest struct {
	Proposal string
	Voter    string
}

type ValidateResponse struct {
	OK bool

	VotingPower     float64
	ValidationError *ValidationError
	VoteStatus      VoteStatus
}

type VoteStatus struct {
	Voted  bool
	Choice json.RawMessage
}

type ValidationError struct {
	Message string
	Code    uint32
}

type PrepareRequest struct {
	Voter    string
	Proposal string
	Choice   json.RawMessage
	Reason   *string
}

type PrepareResponse struct {
	ID        string
	TypedData string
}

type VoteRequest struct {
	ID  string
	Sig string
}

type VoteResponse struct {
	ID      string
	IPFS    string
	Relayer Relayer
}

type Relayer struct {
	Address string
	Receipt string
}
