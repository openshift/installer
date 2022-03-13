//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ovirtsdk

// Service is the interface of all type services.
type Service interface {
	Connection() *Connection
	Path() string
}

// BaseService represents the base for all the services of the SDK. It contains the
// utility methods used by all of them.
type BaseService struct {
	connection *Connection
	path       string
}

func (service *BaseService) Connection() *Connection {
	return service.connection
}

func (service *BaseService) Path() string {
	return service.path
}
