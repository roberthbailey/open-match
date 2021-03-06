// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"open-match.dev/open-match/internal/config"
	internalpb "open-match.dev/open-match/internal/pb"
	"open-match.dev/open-match/internal/rpc"
	"open-match.dev/open-match/internal/statestore"
	"open-match.dev/open-match/pkg/pb"
)

// The service implementing the Backend API that is called to generate matches
// and make assignments for Tickets.
type backendService struct {
	cfg          config.View
	synchronizer internalpb.SynchronizerClient
	store        statestore.Service
	mmfClients   *sync.Map
}

type grpcData struct {
	client pb.MatchFunctionClient
}

type httpData struct {
	client  *http.Client
	baseURL string
}

type mmfResult struct {
	matches []*pb.Match
	err     error
}

var (
	backendServiceLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "app.backend.backend_service",
	})
)

func newBackendService(cfg config.View) (*backendService, error) {
	conn, err := rpc.GRPCClientFromConfig(cfg, "api.synchronizer")
	if err != nil {
		return nil, err
	}

	return &backendService{
		cfg:          cfg,
		synchronizer: internalpb.NewSynchronizerClient(conn),
		store:        statestore.New(cfg),
		mmfClients:   &sync.Map{},
	}, nil
}

// FetchMatches triggers execution of the specfied MatchFunction for each of the
// specified MatchProfiles. Each MatchFunction execution returns a set of
// proposals which are then evaluated to generate results. FetchMatches method
// streams these results back to the caller.
// FetchMatches returns nil unless the context is canceled. FetchMatches moves to the next response candidate if it encounters
// any internal execution failures.
func (s *backendService) FetchMatches(ctx context.Context, req *pb.FetchMatchesRequest) (*pb.FetchMatchesResponse, error) {
	if req.GetConfig() == nil {
		return nil, status.Error(codes.InvalidArgument, ".config is required")
	}
	if req.GetProfile() == nil {
		return nil, status.Error(codes.InvalidArgument, ".profile is required")
	}

	resultChan := make(chan mmfResult, len(req.GetProfile()))

	var syncID string
	if s.synchronizerEnabled() {
		resp, err := s.synchronizer.Register(ctx, &internalpb.RegisterRequest{})
		if err != nil {
			return nil, err
		}

		syncID = resp.GetId()
	}

	err := doFetchMatchesInChannel(ctx, s.cfg, s.mmfClients, req, resultChan)
	if err != nil {
		return nil, err
	}

	proposals, err := doFetchMatchesFilterChannel(ctx, resultChan, len(req.GetProfile()))
	if err != nil {
		return nil, err
	}

	results := proposals
	if s.synchronizerEnabled() {
		resp, err := s.synchronizer.EvaluateProposals(ctx, &internalpb.EvaluateProposalsRequest{
			Id:    syncID,
			Match: proposals})
		if err != nil {
			return nil, err
		}

		results = resp.Match
	}

	return &pb.FetchMatchesResponse{Match: results}, nil
}

func doFetchMatchesInChannel(ctx context.Context, cfg config.View, mmfClients *sync.Map, req *pb.FetchMatchesRequest, resultChan chan<- mmfResult) error {
	var grpcClient pb.MatchFunctionClient
	var httpClient *http.Client
	var baseURL string
	var err error

	configType := req.GetConfig().GetType()
	address := fmt.Sprintf("%s:%d", req.GetConfig().GetHost(), req.GetConfig().GetPort())

	switch configType {
	// MatchFunction Hosted as a GRPC service
	case pb.FunctionConfig_GRPC:
		grpcClient, err = getGRPCClient(cfg, mmfClients, address)
		if err != nil {
			backendServiceLogger.WithFields(logrus.Fields{
				"error":    err.Error(),
				"function": req.GetConfig(),
			}).Error("failed to establish grpc client connection to match function")
			return status.Error(codes.InvalidArgument, "failed to connect to match function")
		}
	// MatchFunction Hosted as a REST service
	case pb.FunctionConfig_REST:
		httpClient, baseURL, err = getHTTPClient(cfg, mmfClients, address)
		if err != nil {
			backendServiceLogger.WithFields(logrus.Fields{
				"error":    err.Error(),
				"function": req.GetConfig(),
			}).Error("failed to establish rest client connection to match function")
			return status.Error(codes.InvalidArgument, "failed to connect to match function")
		}
	default:
		return status.Error(codes.InvalidArgument, "provided match function type is not supported")
	}

	for _, profile := range req.GetProfile() {
		go func(profile *pb.MatchProfile) {
			// Get the match results that will be sent.
			// TODO: The matches returned by the MatchFunction will be sent to the
			// Evaluator to select results. Until the evaluator is implemented,
			// we channel all matches as accepted results.
			switch configType {
			case pb.FunctionConfig_GRPC:
				matches, err := matchesFromGRPCMMF(ctx, profile, grpcClient)
				resultChan <- mmfResult{matches, err}
			case pb.FunctionConfig_REST:
				matches, err := matchesFromHTTPMMF(ctx, profile, httpClient, baseURL)
				resultChan <- mmfResult{matches, err}
			}
		}(profile)
	}

	return nil
}

func doFetchMatchesFilterChannel(ctx context.Context, resultChan <-chan mmfResult, channelSize int) ([]*pb.Match, error) {
	proposals := []*pb.Match{}
	for i := 0; i < channelSize; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-resultChan:
			if result.err != nil {
				return nil, result.err
			}
			proposals = append(proposals, result.matches...)
		}
	}
	return proposals, nil
}

func getHTTPClient(cfg config.View, mmfClients *sync.Map, addr string) (*http.Client, string, error) {
	val, exists := mmfClients.Load(addr)
	data, ok := val.(httpData)
	if !ok || !exists {
		client, baseURL, err := rpc.HTTPClientFromEndpoint(cfg, addr)
		if err != nil {
			return nil, "", err
		}
		data = httpData{client, baseURL}
		mmfClients.Store(addr, data)
	}
	return data.client, data.baseURL, nil
}

func getGRPCClient(cfg config.View, mmfClients *sync.Map, addr string) (pb.MatchFunctionClient, error) {
	val, exists := mmfClients.Load(addr)
	data, ok := val.(grpcData)
	if !ok || !exists {
		conn, err := rpc.GRPCClientFromEndpoint(cfg, addr)
		if err != nil {
			return nil, err
		}
		data = grpcData{pb.NewMatchFunctionClient(conn)}
		mmfClients.Store(addr, data)
	}

	return data.client, nil
}

func matchesFromHTTPMMF(ctx context.Context, profile *pb.MatchProfile, client *http.Client, baseURL string) ([]*pb.Match, error) {
	jsonProfile, err := json.Marshal(profile)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to marshal profile pb to string for profile %s: %s", profile.GetName(), err.Error())
	}

	reqBody, err := json.Marshal(map[string]json.RawMessage{"profile": jsonProfile})
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to marshal request body for profile %s: %s", profile.GetName(), err.Error())
	}

	req, err := http.NewRequest("POST", baseURL+"/v1/matchfunction:run", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to create mmf http request for profile %s: %s", profile.GetName(), err.Error())
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get response from mmf run for proile %s: %s", profile.Name, err.Error())
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			backendServiceLogger.WithError(err).Warning("failed to close response body read closer")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to read from response body for profile %s: %s", profile.Name, err.Error())
	}

	pbResp := &pb.RunResponse{}
	err = json.Unmarshal(body, pbResp)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to unmarshal response body to response pb for profile %s: %s", profile.Name, err.Error())
	}

	return pbResp.GetProposal(), nil
}

// matchesFromGRPCMMF triggers execution of MMFs to fetch match results for each profile.
// These proposals are then sent to evaluator and the results are streamed back on the channel
// that this function returns to the caller.
func matchesFromGRPCMMF(ctx context.Context, profile *pb.MatchProfile, client pb.MatchFunctionClient) ([]*pb.Match, error) {
	// TODO: This code calls user code and could hang. We need to add a deadline here
	// and timeout gracefully to ensure that the ListMatches completes.
	resp, err := client.Run(ctx, &pb.RunRequest{Profile: profile})
	if err != nil {
		backendServiceLogger.WithError(err).Error("failed to run match function for profile")
		return nil, err
	}

	return resp.GetProposal(), nil
}

// AssignTickets sets the specified Assignment on the Tickets for the Ticket
// ids passed.
func (s *backendService) AssignTickets(ctx context.Context, req *pb.AssignTicketsRequest) (*pb.AssignTicketsResponse, error) {
	err := doAssignTickets(ctx, req, s.store)
	if err != nil {
		backendServiceLogger.WithError(err).Error("failed to update assignments for requested tickets")
		return nil, err
	}

	return &pb.AssignTicketsResponse{}, nil
}

func doAssignTickets(ctx context.Context, req *pb.AssignTicketsRequest, store statestore.Service) error {
	err := store.UpdateAssignments(ctx, req.GetTicketId(), req.GetAssignment())
	if err != nil {
		backendServiceLogger.WithError(err).Error("failed to update assignments")
		return err
	}
	for _, id := range req.GetTicketId() {
		err = store.DeindexTicket(ctx, id)
		// Try to deindex all input tickets. Log without returning an error if the deindexing operation failed.
		// TODO: consider retry the index operation
		if err != nil {
			backendServiceLogger.WithError(err).Errorf("failed to deindex ticket %s after updating the assignments", id)
		}
	}

	return nil
}

func (s *backendService) synchronizerEnabled() bool {
	if !s.cfg.IsSet("synchronizer.enabled") {
		return false
	}

	return s.cfg.GetBool("synchronizer.enabled")
}
