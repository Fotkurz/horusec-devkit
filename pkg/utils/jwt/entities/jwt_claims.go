// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package entities

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/Fotkurz/horusec-devkit/pkg/enums/ozzovalidation"
)

type JWTClaims struct {
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func (j *JWTClaims) Validate() error {
	return validation.ValidateStruct(j,
		validation.Field(&j.Username, validation.Required,
			validation.Length(ozzovalidation.Length1, ozzovalidation.Length255)),
		validation.Field(&j.Email, validation.Required,
			validation.Length(ozzovalidation.Length1, ozzovalidation.Length255)),
		validation.Field(&j.Subject, validation.Required, is.UUID, validation.NotIn(uuid.Nil)),
	)
}
