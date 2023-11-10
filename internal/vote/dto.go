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
	TypedData string
}

type VoteRequest struct {
	Voter    string
	Proposal string
	Choice   json.RawMessage
	Reason   *string
	Sig      string
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
