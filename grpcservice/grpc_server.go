/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpcservice

import (
	"context"
	"fmt"
	"net"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/ca17/teamsacs/models"
)

type server struct {
	manager *models.ModelManager
}

func (s *server) GetConfig(ctx context.Context, request *GetConfigRequest) (*GetConfigReply, error) {
	value := s.manager.GetConfigManager().GetConfigValue(request.Type, request.Name)
	return &GetConfigReply{Code: 0, Message: "Success", Value: value}, nil
}

func StartGrpcServer(manager *models.ModelManager) error {
	appconfig := manager.Config
	certfile := path.Join(appconfig.GetPrivateDir(), "teamsacs-grpc.tls.crt")
	keyfile := path.Join(appconfig.GetPrivateDir(), "teamsacs-grpc.tls.key")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", appconfig.Grpc.Host, appconfig.Grpc.Port))
	if err != nil {
		return err
	}
	creds, err := credentials.NewServerTLSFromFile(certfile, keyfile)
	if err != nil {
		return err
	}
	s := grpc.NewServer(grpc.Creds(creds))
	RegisterTeamsacsServiceServer(s, &server{manager: manager})
	reflection.Register(s)
	return s.Serve(lis)
}
