// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package proto

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	mockUtils "github.com/Fotkurz/horusec-devkit/pkg/utils/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) IsAuthorized(
	_ context.Context, _ *IsAuthorizedData, _ ...grpc.CallOption) (*IsAuthorizedResponse, error) {
	args := m.MethodCalled("IsAuthorized")

	return args.Get(0).(*IsAuthorizedResponse), mockUtils.ReturnNilOrError(args, 1)
}

func (m *Mock) GetAccountInfo(_ context.Context, _ *GetAccountData,
	_ ...grpc.CallOption) (*GetAccountDataResponse, error) {
	args := m.MethodCalled("GetAccountInfo")

	return args.Get(0).(*GetAccountDataResponse), mockUtils.ReturnNilOrError(args, 1)
}

func (m *Mock) GetAuthConfig(_ context.Context, _ *GetAuthConfigData,
	_ ...grpc.CallOption) (*GetAuthConfigResponse, error) {
	args := m.MethodCalled("GetAuthConfig")

	return args.Get(0).(*GetAuthConfigResponse), mockUtils.ReturnNilOrError(args, 1)
}
