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

package linux

import (
	"gvisor.dev/gvisor/pkg/marshal"
)

// Constants for error origin.
const (
	SO_EE_ORIGIN_NONE  = 0
	SO_EE_ORIGIN_LOCAL = 1
	SO_EE_ORIGIN_ICMP  = 2
	SO_EE_ORIGIN_ICMP6 = 3
)

// SockExtendedErr represents struct sock_extended_err in Linux defined in
// include/uapi/linux/errqueue.h.
//
// +marshal
type SockExtendedErr struct {
	Errno  uint32
	Origin uint8
	Type   uint8
	Code   uint8
	Pad    uint8
	Info   uint32
	Data   uint32
}

// IsICMPErr implements SockErrCMsg.IsICMPErr and is equivalent to
// net/core/skbuff.c:is_icmp_err_skb().
func (ee *SockExtendedErr) IsICMPErr() bool {
	if ee == nil {
		return false
	}

	return ee.Origin == SO_EE_ORIGIN_ICMP || ee.Origin == SO_EE_ORIGIN_ICMP6
}

// SockErrCMsg represents the IP*_RECVERR control message.
type SockErrCMsg interface {
	marshal.Marshallable

	CMsgLevel() uint32
	CMsgType() uint32
	IsICMPErr() bool
}

// SockErrCMsgIPv4 is the IP_RECVERR control message used in
// recvmsg(MSG_ERRQUEUE) by ipv4 sockets. This is equilavent to `struct errhdr`
// defined in net/ipv4/ip_sockglue.c:ip_recv_error().
//
// +marshal
type SockErrCMsgIPv4 struct {
	SockExtendedErr
	Offender SockAddrInet
}

var _ SockErrCMsg = (*SockErrCMsgIPv4)(nil)

// CMsgLevel implements SockErrCMsg.CMsgLevel.
func (*SockErrCMsgIPv4) CMsgLevel() uint32 {
	return SOL_IP
}

// CMsgType implements SockErrCMsg.CMsgType.
func (*SockErrCMsgIPv4) CMsgType() uint32 {
	return IP_RECVERR
}

// SockErrCMsgIPv6 is the IPV6_RECVERR control message used in
// recvmsg(MSG_ERRQUEUE) by ipv6 sockets. This is equilavent to `struct errhdr`
// defined in net/ipv6/datagram.c:ipv6_recv_error().
//
// +marshal
type SockErrCMsgIPv6 struct {
	SockExtendedErr
	Offender SockAddrInet6
}

var _ SockErrCMsg = (*SockErrCMsgIPv6)(nil)

// CMsgLevel implements SockErrCMsg.CMsgLevel.
func (*SockErrCMsgIPv6) CMsgLevel() uint32 {
	return SOL_IPV6
}

// CMsgType implements SockErrCMsg.CMsgType.
func (*SockErrCMsgIPv6) CMsgType() uint32 {
	return IPV6_RECVERR
}
