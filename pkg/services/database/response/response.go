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

type IResponse interface {
	GetRowsAffected() int
	SetRowsAffected(value int)
	GetData() interface{}
	SetData(value interface{})
	GetError() error
	SetError(value error)
}

type Response struct {
	err          error
	rowsAffected int
	data         interface{}
}

func NewResponse(rowsAffected int, err error, data interface{}) IResponse {
	return &Response{
		err:          err,
		rowsAffected: rowsAffected,
		data:         data,
	}
}

func (r *Response) GetRowsAffected() int {
	return r.rowsAffected
}

func (r *Response) SetRowsAffected(value int) {
	r.rowsAffected = value
}

func (r *Response) GetData() interface{} {
	return r.data
}

func (r *Response) SetData(value interface{}) {
	r.data = value
}

func (r *Response) GetError() error {
	return r.err
}

func (r *Response) SetError(value error) {
	r.err = value
}
