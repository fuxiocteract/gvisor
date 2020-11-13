// Copyright 2020 The gVisor Authors.
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

package tcp_zero_receive_window_test

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/test/packetimpact/testbench"
)

func init() {
	testbench.RegisterFlags(flag.CommandLine)
}

func TestZeroReceiveWindow(t *testing.T) {
	for _, payloadLen := range []int{64, 512, 1024} {
		t.Run(fmt.Sprintf("TestZeroReceiveWindow_with_%dbytes_payload", payloadLen), func(t *testing.T) {
			dut := testbench.NewDUT(t)
			defer dut.TearDown()
			listenFd, remotePort := dut.CreateListener(t, unix.SOCK_STREAM, unix.IPPROTO_TCP, 1)
			defer dut.Close(t, listenFd)
			conn := testbench.NewTCPIPv4(t, testbench.TCP{DstPort: &remotePort}, testbench.TCP{SrcPort: &remotePort})
			defer conn.Close(t)

			conn.Connect(t)
			acceptFd, _ := dut.Accept(t, listenFd)
			defer dut.Close(t, acceptFd)

			dut.SetSockOptInt(t, acceptFd, unix.IPPROTO_TCP, unix.TCP_NODELAY, 1)

			samplePayload := &testbench.Payload{Bytes: make([]byte, payloadLen)}

			// Expect the DUT to eventually advertize zero receive window.
			// The test would timeout otherwise.
			for {
				conn.Send(t, testbench.TCP{Flags: testbench.Uint8(header.TCPFlagAck | header.TCPFlagPsh)}, samplePayload)
				gotTCP, err := conn.Expect(t, testbench.TCP{Flags: testbench.Uint8(header.TCPFlagAck)}, time.Second)
				if err != nil {
					t.Fatalf("expected packet was not received: %s", err)
				}
				if *gotTCP.WindowSize == 0 {
					break
				}
			}
		})
	}
}
