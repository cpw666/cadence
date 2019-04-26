// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cassandra

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/uber/cadence/environment"
)

type (
	HandlerTestSuite struct {
		*require.Assertions // override suite.Suite.Assertions with require.Assertions; this means that s.NotNil(nil) will stop the test, not merely log an error
		suite.Suite
	}
)

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.Assertions = require.New(s.T()) // Have to define our overridden assertions in the test setup. If we did it earlier, s.T() will return nil
}

func (s *HandlerTestSuite) TestValidateSetupSchemaConfig() {

	config := new(SetupSchemaConfig)
	s.assertValidateSetupFails(config)

	config.CassHosts = environment.GetCassandraAddress()
	s.assertValidateSetupFails(config)

	config.CassKeyspace = "test-keyspace"
	s.assertValidateSetupFails(config)

	config.InitialVersion = "0.1"
	config.DisableVersioning = true
	config.SchemaFilePath = ""
	s.assertValidateSetupFails(config)

	config.InitialVersion = "0.1"
	config.DisableVersioning = true
	config.SchemaFilePath = "/tmp/foo.cql"
	s.assertValidateSetupFails(config)

	config.InitialVersion = ""
	config.DisableVersioning = true
	config.SchemaFilePath = ""
	s.assertValidateSetupFails(config)

	config.InitialVersion = "0.1"
	config.DisableVersioning = false
	config.SchemaFilePath = "/tmp/foo.cql"
	s.assertValidateSetupSucceeds(config)

	config.InitialVersion = "0.1"
	config.DisableVersioning = false
	config.SchemaFilePath = ""
	s.assertValidateSetupSucceeds(config)

	config.InitialVersion = ""
	config.DisableVersioning = true
	config.SchemaFilePath = "/tmp/foo.cql"
	s.assertValidateSetupSucceeds(config)
}

func (s *HandlerTestSuite) TestValidateUpdateSchemaConfig() {

	config := new(UpdateSchemaConfig)
	s.assertValidateUpdateFails(config)

	config.CassHosts = environment.GetCassandraAddress()
	s.assertValidateUpdateFails(config)

	config.CassKeyspace = "test-keyspace"
	s.assertValidateUpdateFails(config)

	config.SchemaDir = "/tmp"
	config.TargetVersion = "abc"
	s.assertValidateUpdateFails(config)

	config.SchemaDir = "/tmp"
	config.TargetVersion = ""
	s.assertValidateUpdateSucceeds(config)

	config.SchemaDir = "/tmp"
	config.TargetVersion = "1.2"
	s.assertValidateUpdateSucceeds(config)

	config.SchemaDir = "/tmp"
	config.TargetVersion = "v1.2"
	s.assertValidateUpdateSucceeds(config)
	s.Equal("1.2", config.TargetVersion)
}

func (s *HandlerTestSuite) TestValidateCreateKeyspaceConfig() {
	config := new(CreateKeyspaceConfig)
	s.NotNil(validateCreateKeyspaceConfig(config))
	config.CassHosts = environment.GetCassandraAddress()
	s.NotNil(validateCreateKeyspaceConfig(config))
	config.CassKeyspace = "foobar"
	s.Nil(validateCreateKeyspaceConfig(config))
}

func (s *HandlerTestSuite) assertValidateSetupSucceeds(input *SetupSchemaConfig) {
	err := validateSetupSchemaConfig(input)
	s.Nil(err)
}

func (s *HandlerTestSuite) assertValidateSetupFails(input *SetupSchemaConfig) {
	err := validateSetupSchemaConfig(input)
	s.NotNil(err)
	_, ok := err.(*ConfigError)
	s.True(ok)
}

func (s *HandlerTestSuite) assertValidateUpdateSucceeds(input *UpdateSchemaConfig) {
	err := validateUpdateSchemaConfig(input)
	s.Nil(err)
}

func (s *HandlerTestSuite) assertValidateUpdateFails(input *UpdateSchemaConfig) {
	err := validateUpdateSchemaConfig(input)
	s.NotNil(err)
	_, ok := err.(*ConfigError)
	s.True(ok)
}