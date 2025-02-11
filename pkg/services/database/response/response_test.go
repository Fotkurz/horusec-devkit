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

package response

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fotkurz/horusec-devkit/pkg/services/database/enums"
)

func TestNewResponse(t *testing.T) {
	t.Run("should create new response with arguments", func(t *testing.T) {
		databaseResponse := NewResponse(10, errors.New("test"), "data")

		assert.Equal(t, 10, databaseResponse.GetRowsAffected())
		assert.Equal(t, errors.New("test"), databaseResponse.GetError())
		assert.Equal(t, "data", databaseResponse.GetData())
	})
}

func TestGetRowsAffected(t *testing.T) {
	t.Run("should success get rows affected", func(t *testing.T) {
		databaseResponse := &Response{rowsAffected: 5}

		assert.Equal(t, 5, databaseResponse.GetRowsAffected())
	})
}

func TestGetError(t *testing.T) {
	t.Run("should success get error", func(t *testing.T) {
		databaseResponse := &Response{err: errors.New("error")}

		assert.Equal(t, errors.New("error"), databaseResponse.GetError())
	})
}

func TestGetData(t *testing.T) {
	t.Run("should success get data", func(t *testing.T) {
		databaseResponse := &Response{data: "test"}

		assert.Equal(t, "test", databaseResponse.GetData())
	})
}

func TestGetErrorExceptNotFound(t *testing.T) {
	t.Run("should return nil when not found error", func(t *testing.T) {
		databaseResponse := &Response{err: enums.ErrorNotFoundRecords}

		assert.NoError(t, databaseResponse.GetErrorExceptNotFound())
	})

	t.Run("should return error when it is something different than not found", func(t *testing.T) {
		databaseResponse := &Response{err: errors.New("test")}

		assert.Error(t, databaseResponse.GetErrorExceptNotFound())
	})
}
