//go:generate protoc --go_out=. --go-grpc_out=. ./base.proto
//go:generate protoc --go_out=. --go-grpc_out=. ./dao.proto
//go:generate protoc --go_out=. --go-grpc_out=. ./proposal.proto
package internalapi
