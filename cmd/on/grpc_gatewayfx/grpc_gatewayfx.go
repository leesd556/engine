/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc_gatewayfx

import (
	"context"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/grpc_gateway/api"
	"github.com/it-chain/engine/grpc_gateway/infra"
	"github.com/it-chain/engine/grpc_gateway/infra/adapter"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewGrpcHostService,
		NewConnectionApi,
		adapter.NewConnectionCommandHandler,
	),
	fx.Invoke(
		RegisterHandlers,
		InitgRPCServer,
	),
)

func NewGrpcHostService(conf *conf.Configuration, publisher *pubsub.TopicPublisher) *infra.GrpcHostService {
	priKey, pubKey := infra.LoadKeyPair(conf.Engine.KeyPath, "ECDSA256")
	hostService := infra.NewGrpcHostService(priKey, pubKey, publisher.Publish)
	return hostService
}

func NewConnectionApi(hostService *infra.GrpcHostService, eventService common.EventService) *api.ConnectionApi {
	return api.NewConnectionApi(hostService, eventService)
}

func RegisterHandlers(connectionCommandHandler *adapter.ConnectionCommandHandler, server *rpc.Server) {
	logger.Infof(nil, "[Main] gRPC-Gateway is starting")
	if err := server.Register("connection.create", connectionCommandHandler.HandleCreateConnectionCommand); err != nil {
		panic(err)
	}

	if err := server.Register("connection.list", connectionCommandHandler.HandleGetConnectionListCommand); err != nil {
		panic(err)
	}

	if err := server.Register("connection.close", connectionCommandHandler.HandleCloseConnectionCommand); err != nil {
		panic(err)
	}
}

func InitgRPCServer(lifecycle fx.Lifecycle, config *conf.Configuration, hostService *infra.GrpcHostService) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			go hostService.Listen(config.GrpcGateway.Address + ":" + config.GrpcGateway.Port)
			return nil
		},
		OnStop: func(context context.Context) error {
			connections, _ := hostService.GetAllConnections()
			for _, connection := range connections {
				hostService.CloseConnection(connection.ConnectionId)
			}
			hostService.Stop()
			return nil
		},
	})
}
