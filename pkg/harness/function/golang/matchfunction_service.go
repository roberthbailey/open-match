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

// Package golang provides the Match Making Function service for Open Match golang harness.
package golang

import (
	"context"
	"io"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"open-match.dev/open-match/internal/config"
	"open-match.dev/open-match/internal/rpc"
	"open-match.dev/open-match/pkg/pb"

	"github.com/sirupsen/logrus"
)

var (
	matchfunctionLogger = logrus.WithFields(logrus.Fields{
		"app":       "openmatch",
		"component": "harness.golang.matchfunction_service",
	})
)

// MatchFunction is the function signature for the Match Making Function (MMF) to be implemented by the user.
// The harness will pass the Rosters and PlayerPool for the match profile to this
// function and it will return the Rosters to be populated in the proposal.
// Input:
//  - MatchFunctionParams:
//			A structure that defines the resources that are available to the match function.
//			Developers can choose to add context to the structure such that match function has the ability
//			to cancel a stream response/request or to limit match function by sharing a static and protected view.
type MatchFunction func(*MatchFunctionParams) []*pb.Match

// matchFunctionService implements pb.MatchFunctionServer, the server generated
// by compiling the protobuf, by fulfilling the pb.MatchFunctionServer interface.
type matchFunctionService struct {
	cfg           config.View
	function      MatchFunction
	mmlogicClient pb.MmLogicClient
}

// MatchFunctionParams : This harness example defines a protected view for the match function.
//  - logger:
//			A logger used to generate error/debug logs
//  - profileName
//			The name of profile
//  - properties:
//			'Properties' of a MatchObject.
//  - rosters:
//			An array of Rosters. By convention, your input Roster contains players already in
//          the match, and the names of pools to search when trying to fill an empty slot.
//  - poolNameToTickets:
//			A map that contains mappings from pool name to a list of tickets that satisfied the filters in the pool
type MatchFunctionParams struct {
	Logger            *logrus.Entry
	ProfileName       string
	Properties        *structpb.Struct
	Rosters           []*pb.Roster
	PoolNameToTickets map[string][]*pb.Ticket
}

// Run is this harness's implementation of the gRPC call defined in api/matchfunction.proto.
func (s *matchFunctionService) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {
	poolNameToTickets, err := s.getMatchManifest(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// The matchfunction takes in some half-filled/empty rosters, a property bag, and a map[poolNames]tickets to generate match proposals
	mfView := &MatchFunctionParams{
		Logger:            matchfunctionLogger,
		ProfileName:       req.GetProfile().GetName(),
		Properties:        req.GetProfile().GetProperties(),
		Rosters:           req.GetProfile().GetRoster(),
		PoolNameToTickets: poolNameToTickets,
	}
	// Run the customize match function!
	proposals := s.function(mfView)
	matchfunctionLogger.WithFields(logrus.Fields{
		"proposals": proposals,
	}).Trace("proposals returned by match function")
	return &pb.RunResponse{Proposal: proposals}, nil
}

func newMatchFunctionService(cfg config.View, fs *FunctionSettings) (*matchFunctionService, error) {
	conn, err := rpc.GRPCClientFromConfig(cfg, "api.mmlogic")
	if err != nil {
		harnessLogger.Errorf("Failed to get MMLogic connection, %v.", err)
		return nil, err
	}

	mmfService := &matchFunctionService{cfg: cfg, function: fs.Func, mmlogicClient: pb.NewMmLogicClient(conn)}
	return mmfService, nil
}

// getMatchManifest fetches all the data needed from the mmlogic API.
func (s *matchFunctionService) getMatchManifest(ctx context.Context, req *pb.RunRequest) (map[string][]*pb.Ticket, error) {
	poolNameToTickets := make(map[string][]*pb.Ticket)
	filterPools := req.GetProfile().GetPool()

	for _, pool := range filterPools {
		qtClient, err := s.mmlogicClient.QueryTickets(ctx, &pb.QueryTicketsRequest{Pool: pool})
		if err != nil {
			matchfunctionLogger.WithError(err).Error("Failed to get queryTicketClient from mmlogic.")
			return nil, err
		}

		// Aggregate tickets by poolName
		poolTickets := make([]*pb.Ticket, 0)
		for {
			qtResponse, err := qtClient.Recv()
			if err == io.EOF {
				matchfunctionLogger.Trace("Received all results from the queryTicketClient.")
				// Break when all results are received
				break
			}

			if err != nil {
				matchfunctionLogger.WithError(err).Error("Failed to receive a response from the queryTicketClient.")
				return nil, err
			}
			poolTickets = append(poolTickets, qtResponse.Ticket...)
		}
		poolNameToTickets[pool.GetName()] = poolTickets
	}

	return poolNameToTickets, nil
}
