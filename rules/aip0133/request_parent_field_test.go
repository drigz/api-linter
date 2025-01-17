// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aip0133

import (
	"testing"

	"github.com/googleapis/api-linter/rules/internal/testutils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func TestParentField(t *testing.T) {
	// Set up the testing permutations.
	tests := []struct {
		testName     string
		messageName  string
		messageField *field
		problems     testutils.Problems
		problemDesc  func(m *desc.MessageDescriptor) desc.Descriptor
	}{
		{
			"Valid",
			"CreateBookRequest",
			&field{"parent", builder.FieldTypeString()},
			testutils.Problems{},
			nil},
		{
			"MissingField",
			"CreateBookRequest",
			nil,
			testutils.Problems{{Message: "parent"}},
			nil,
		},
		{
			"InvalidType",
			"CreateBookRequest",
			&field{"parent", builder.FieldTypeDouble()},
			testutils.Problems{{Suggestion: "string"}},
			func(m *desc.MessageDescriptor) desc.Descriptor {
				return m.FindFieldByName("parent")
			},
		},
	}

	// Run each test individually.
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create an appropriate message descriptor.
			messageBuilder := builder.NewMessage(test.messageName)

			if test.messageField != nil {
				messageBuilder.AddField(
					builder.NewField(test.messageField.fieldName, test.messageField.fieldType),
				)
			}

			message, err := messageBuilder.Build()
			if err != nil {
				t.Fatalf("Could not build %s message.", test.messageName)
			}

			// Determine the descriptor that a failing test will attach to.
			var problemDesc desc.Descriptor = message
			if test.problemDesc != nil {
				problemDesc = test.problemDesc(message)
			}

			// Run the lint rule, and establish that it returns the correct problems.
			problems := parentField.Lint(message.GetFile())
			if diff := test.problems.SetDescriptor(problemDesc).Diff(problems); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
