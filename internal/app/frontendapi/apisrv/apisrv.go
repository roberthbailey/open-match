/*
package apisrv provides an implementation of the gRPC server defined in ../../../api/protobuf-spec/frontend.proto.

Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package apisrv

import (
	"context"
	"errors"
	"net"

	"github.com/GoogleCloudPlatform/open-match/internal/expbo"
	"github.com/GoogleCloudPlatform/open-match/internal/metrics"
	frontend "github.com/GoogleCloudPlatform/open-match/internal/pb"
	redishelpers "github.com/GoogleCloudPlatform/open-match/internal/statestorage/redis"
	"github.com/GoogleCloudPlatform/open-match/internal/statestorage/redis/playerindices"
	"github.com/GoogleCloudPlatform/open-match/internal/statestorage/redis/redispb"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

// Logrus structured logging setup
var (
	feLogFields = log.Fields{
		"app":       "openmatch",
		"component": "frontend",
	}
	feLog = log.WithFields(feLogFields)
)

// FrontendAPI implements frontend.ApiServer, the server generated by compiling
// the protobuf, by fulfilling the frontend.APIClient interface.
type FrontendAPI struct {
	grpc *grpc.Server
	cfg  *viper.Viper
	pool *redis.Pool
}
type frontendAPI FrontendAPI

// New returns an instantiated srvice
func New(cfg *viper.Viper, pool *redis.Pool) *FrontendAPI {
	s := FrontendAPI{
		pool: pool,
		grpc: grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{})),
		cfg:  cfg,
	}

	// Add a hook to the logger to auto-count log lines for metrics output thru OpenCensus
	log.AddHook(metrics.NewHook(FeLogLines, KeySeverity))

	// Register gRPC server
	frontend.RegisterFrontendServer(s.grpc, (*frontendAPI)(&s))
	feLog.Info("Successfully registered gRPC server")
	return &s
}

// Open starts the api grpc service listening on the configured port.
func (s *FrontendAPI) Open() error {
	ln, err := net.Listen("tcp", ":"+s.cfg.GetString("api.frontend.port"))
	if err != nil {
		feLog.WithFields(log.Fields{
			"error": err.Error(),
			"port":  s.cfg.GetInt("api.frontend.port"),
		}).Error("net.Listen() error")
		return err
	}
	feLog.WithFields(log.Fields{"port": s.cfg.GetInt("api.frontend.port")}).Info("TCP net listener initialized")

	go func() {
		err := s.grpc.Serve(ln)
		if err != nil {
			feLog.WithFields(log.Fields{"error": err.Error()}).Error("gRPC serve() error")
		}
		feLog.Info("serving gRPC endpoints")
	}()

	return nil
}

// CreatePlayer is this service's implementation of the CreatePlayer gRPC method defined in frontend.proto
func (s *frontendAPI) CreatePlayer(ctx context.Context, group *frontend.Player) (*frontend.Result, error) {

	// Create context for tagging OpenCensus metrics.
	funcName := "CreatePlayer"
	fnCtx, _ := tag.New(ctx, tag.Insert(KeyMethod, funcName))

	// Write group
	err := redispb.MarshalToRedis(ctx, s.pool, group, s.cfg.GetInt("redis.expirations.player"))
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}

	// Index group
	err = playerindices.Create(ctx, s.pool, s.cfg, *group)
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}

	// Return success.
	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.Result{Success: true, Error: ""}, err

}

// DeletePlayer is this service's implementation of the DeletePlayer gRPC method defined in frontend.proto
func (s *frontendAPI) DeletePlayer(ctx context.Context, group *frontend.Player) (*frontend.Result, error) {

	// Create context for tagging OpenCensus metrics.
	funcName := "DeletePlayer"
	fnCtx, _ := tag.New(ctx, tag.Insert(KeyMethod, funcName))

	// Deindex this player; at that point they don't show up in MMFs anymore.  We can then delete
	// their actual player object from Redis later.
	err := playerindices.Delete(ctx, s.pool, s.cfg, group.Id)
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Error("State storage error")

		stats.Record(fnCtx, FeGrpcErrors.M(1))
		return &frontend.Result{Success: false, Error: err.Error()}, err
	}
	// Kick off delete but don't wait for it to complete.
	go s.deletePlayer(group.Id)

	stats.Record(fnCtx, FeGrpcRequests.M(1))
	return &frontend.Result{Success: true, Error: ""}, err

}

// deletePlayer is a 'lazy' player delete
// It should always be called as a goroutine and should only be called after
// confirmation that a player has been deindexed (and therefore MMF's can't
// find the player to read them anyway)
// As a final action, it also kicks off a lazy delete of the player's metadata
func (s *frontendAPI) deletePlayer(id string) {

	err := redishelpers.Delete(context.Background(), s.pool, id)
	if err != nil {
		feLog.WithFields(log.Fields{
			"error":     err.Error(),
			"component": "statestorage",
		}).Warn("Error deleting player from state storage, this could leak state storage memory but is usually not a fatal error")
	}
	go playerindices.DeleteMeta(context.Background(), s.pool, id)
}

// GetUpdates is this service's implementation of the GetUpdates gRPC method defined in frontend.proto
func (s *frontendAPI) GetUpdates(p *frontend.Player, assignmentStream frontend.Frontend_GetUpdatesServer) error {
	// Get cancellable context
	ctx, cancel := context.WithCancel(assignmentStream.Context())
	defer cancel()

	// Create context for tagging OpenCensus metrics.
	funcName := "GetAssignment"
	fnCtx, _ := tag.New(ctx, tag.Insert(KeyMethod, funcName))

	watcherBO := backoff.NewExponentialBackOff()
	if err := expbo.UnmarshalExponentialBackOff(s.cfg.GetString("api.frontend.backoff"), watcherBO); err != nil {
		feLog.WithError(err).Warn("Could not parse backoff string, using default backoff parameters for Player watcher")
	}

	// We have to stop Watcher manually because in a normal case client closes channel before the timeout
	watcherCtx, stopWatcher := context.WithCancel(context.Background())
	defer stopWatcher()
	watcherBOCtx := backoff.WithContext(watcherBO, watcherCtx)

	// get and return connection string
	watchChan := redispb.PlayerWatcher(watcherBOCtx, s.pool, *p) // watcher() runs the appropriate Redis commands.

	for {
		select {
		case <-ctx.Done():
			// Context cancelled
			feLog.WithField("playerid", p.Id).Info("client closed connection successfully")
			stats.Record(fnCtx, FeGrpcRequests.M(1))
			return nil

		case a, ok := <-watchChan:
			if !ok {
				// Timeout reached without client closing connection
				err := errors.New("server timeout reached without client closing connection")
				feLog.WithFields(log.Fields{
					"error":     err.Error(),
					"component": "statestorage",
					"playerid":  p.Id,
				}).Error("State storage error")

				// Count errors for metrics
				errTag, _ := tag.NewKey("errtype")
				fnCtx, _ := tag.New(ctx, tag.Insert(errTag, "watch_timeout"))
				stats.Record(fnCtx, FeGrpcErrors.M(1))
				//TODO: we could generate a frontend.player message with an error
				//field and stream it to the client before throwing the error here
				//if we wanted to send more useful client retry information
				return err
			}

			feLog.WithFields(log.Fields{
				"assignment": a.Assignment,
				"playerid":   a.Id,
				"status":     a.Status,
				"error":      a.Error,
			}).Info("updating client")
			assignmentStream.Send(&a)
			stats.Record(fnCtx, FeGrpcStreamedResponses.M(1))
		}
	}

}
