package http

import (
	"GoFeed/internal/generated/api/go_feed"
	grpc_handle "GoFeed/internal/handler/grpc"
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type accountHandler struct {
	clientPool *grpcClientPool
}

func NewaccountHandler(clientPool *grpcClientPool) *accountHandler {
	return &accountHandler{clientPool: clientPool}
}

func (h accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		AccountName string `json:"account_name"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if body.AccountName == "" {
		WriteError(w, http.StatusBadRequest, "account_name is required")
		return
	}
	if body.Password == "" {
		WriteError(w, http.StatusBadRequest, "password is required")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)

	output, err := client.CreateAccount(ctx, &go_feed.CreateAccountRequest{
		AccountName: body.AccountName,
		Password:    body.Password,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, output)
}

func (h accountHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed, expected POST")
		return
	}

	var body struct {
		AccountName string `json:"account_name"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if body.AccountName == "" {
		WriteError(w, http.StatusBadRequest, "account_name is required and must be a string")
		return
	}

	if body.Password == "" {
		WriteError(w, http.StatusBadRequest, "password is required and must be a string")
		return
	}

	client, err := h.clientPool.GetClient()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mData := metadata.MD{}
	ctx := metadata.NewOutgoingContext(context.Background(), mData)
	var header metadata.MD

	// Call gRPC CreateSession
	output, err := client.CreateSession(ctx, &go_feed.CreateSessionRequest{
		AccountName: body.AccountName,
		Password:    body.Password,
	}, grpc.Header(&header))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create session: "+err.Error())
		return
	}

	// Extract token from gRPC header
	token := header.Get(grpc_handle.AuthTokenMetadataName)
	if len(token) == 0 {
		WriteError(w, http.StatusInternalServerError, "Token not found in gRPC metadata")
		return
	}

	// Write success response
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"accountId": output.AccountId,
		"token":     token[0],
	})
}
