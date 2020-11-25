// Copyright 2018 The gVisor Authors.
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

package header_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gvisor.dev/gvisor/pkg/tcpip/header"
)

func TestOptionsSerializer(t *testing.T) {
	optCases := []struct {
		name   string
		option []header.IPv4SerializableOption
		expect []byte
	}{
		{
			name: "NOP",
			option: []header.IPv4SerializableOption{
				&header.IPv4SerializableNOPOption{},
			},
			expect: []byte{1, 0, 0, 0},
		},
		{
			name: "ListEnd",
			option: []header.IPv4SerializableOption{
				&header.IPv4SerializableListEndOption{},
			},
			expect: []byte{0, 0, 0, 0},
		},
		{
			name: "RouterAlert",
			option: []header.IPv4SerializableOption{
				&header.IPv4SerializableRouterAlertOption{},
			},
			expect: []byte{148, 4, 0, 0},
		}, {
			name: "NOP and RouterAlert",
			option: []header.IPv4SerializableOption{
				&header.IPv4SerializableNOPOption{},
				&header.IPv4SerializableRouterAlertOption{},
			},
			expect: []byte{1, 148, 4, 0, 0, 0, 0, 0},
		},
	}

	for _, opt := range optCases {
		t.Run(opt.name, func(t *testing.T) {
			s := header.IPv4OptionsSerializer(opt.option)
			l := s.Length()
			if got, want := l, len(opt.expect); got != want {
				t.Fatalf("s.Length() = %d, want = %d", got, want)
			}
			b := make([]byte, l, l)
			if serializedLength := s.Serialize(b); serializedLength != l {
				t.Fatalf("s.Serialize(_) = %d, want %d", serializedLength, l)
			}
			if diff := cmp.Diff(b, opt.expect); diff != "" {
				t.Errorf("mismatched serialized option (-want +got):\n%s", diff)
			}
		})
	}
}
