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

package tcpip

import (
	"sync/atomic"

	"gvisor.dev/gvisor/pkg/sync"
)

// SocketOptionsHandler holds methods that help define endpoint specific
// behavior for socket options.
// These must be implemented by endpoints to:
// - Get notified when socket level options are set.
// - Provide endpoint specific socket options.
type SocketOptionsHandler interface {
	// OnReuseAddressSet is invoked when SO_REUSEADDR is set for an endpoint.
	OnReuseAddressSet(v bool)

	// OnReusePortSet is invoked when SO_REUSEPORT is set for an endpoint.
	OnReusePortSet(v bool)

	// OnKeepAliveSet is invoked when SO_KEEPALIVE is set for an endpoint.
	OnKeepAliveSet(v bool)

	// IsListening is invoked to fetch SO_ACCEPTCONN option value for an
	// endpoint. It is used to indicate if the socket is a listening socket.
	IsListening() bool

	// LastError is invoked when SO_ERROR is read for an endpoint..
	LastError() *Error
}

// DefaultSocketOptionsHandler is an embeddable type that implements no-op
// implementations for SocketOptionsHandler methods.
type DefaultSocketOptionsHandler struct{}

var _ SocketOptionsHandler = (*DefaultSocketOptionsHandler)(nil)

// OnReuseAddressSet implements SocketOptionsHandler.OnReuseAddressSet.
func (*DefaultSocketOptionsHandler) OnReuseAddressSet(bool) {}

// OnReusePortSet implements SocketOptionsHandler.OnReusePortSet.
func (*DefaultSocketOptionsHandler) OnReusePortSet(bool) {}

// OnKeepAliveSet implements SocketOptionsHandler.OnKeepAliveSet.
func (*DefaultSocketOptionsHandler) OnKeepAliveSet(bool) {}

// IsListening implements SocketOptionsHandler.IsListening.
func (*DefaultSocketOptionsHandler) IsListening() bool { return false }

// LastError implements SocketOptionsHandler.LastError.
func (*DefaultSocketOptionsHandler) LastError() *Error {
	return nil
}

// SocketOptions contains all the variables which store values for SOL_SOCKET
// level options.
//
// +stateify savable
type SocketOptions struct {
	handler SocketOptionsHandler

	// These fields are accessed and modified using atomic operations.

	// broadcastEnabled determines whether datagram sockets are allowed to
	// send packets to a broadcast address.
	broadcastEnabled uint32

	// passCredEnabled determines whether SCM_CREDENTIALS socket control
	// messages are enabled.
	passCredEnabled uint32

	// noChecksumEnabled determines whether UDP checksum is disabled while
	// transmitting for this socket.
	noChecksumEnabled uint32

	// reuseAddressEnabled determines whether Bind() should allow reuse of
	// local address.
	reuseAddressEnabled uint32

	// reusePortEnabled determines whether to permit multiple sockets to be
	// bound to an identical socket address.
	reusePortEnabled uint32

	// keepAliveEnabled determines whether TCP keepalive is enabled for this
	// socket.
	keepAliveEnabled uint32

	// mu protects the access to the below fields.
	mu sync.Mutex `state:"nosave"`

	// linger determines the amount of time TCP socket should linger before
	// close.
	linger LingerOption

	// bindToDevice determines the device to which the socket is bound.
	bindToDevice NICID
}

// InitHandler initializes the handler. This must be called before using the
// socket options utility.
func (so *SocketOptions) InitHandler(handler SocketOptionsHandler) {
	so.handler = handler
}

func storeAtomicBool(addr *uint32, v bool) {
	var val uint32
	if v {
		val = 1
	}
	atomic.StoreUint32(addr, val)
}

// GetBroadcast gets value for SO_BROADCAST option.
func (so *SocketOptions) GetBroadcast() bool {
	return atomic.LoadUint32(&so.broadcastEnabled) != 0
}

// SetBroadcast sets value for SO_BROADCAST option.
func (so *SocketOptions) SetBroadcast(v bool) {
	storeAtomicBool(&so.broadcastEnabled, v)
}

// GetPassCred gets value for SO_PASSCRED option.
func (so *SocketOptions) GetPassCred() bool {
	return atomic.LoadUint32(&so.passCredEnabled) != 0
}

// SetPassCred sets value for SO_PASSCRED option.
func (so *SocketOptions) SetPassCred(v bool) {
	storeAtomicBool(&so.passCredEnabled, v)
}

// GetNoChecksum gets value for SO_NO_CHECK option.
func (so *SocketOptions) GetNoChecksum() bool {
	return atomic.LoadUint32(&so.noChecksumEnabled) != 0
}

// SetNoChecksum sets value for SO_NO_CHECK option.
func (so *SocketOptions) SetNoChecksum(v bool) {
	storeAtomicBool(&so.noChecksumEnabled, v)
}

// GetReuseAddress gets value for SO_REUSEADDR option.
func (so *SocketOptions) GetReuseAddress() bool {
	return atomic.LoadUint32(&so.reuseAddressEnabled) != 0
}

// SetReuseAddress sets value for SO_REUSEADDR option.
func (so *SocketOptions) SetReuseAddress(v bool) {
	storeAtomicBool(&so.reuseAddressEnabled, v)
	so.handler.OnReuseAddressSet(v)
}

// GetReusePort gets value for SO_REUSEPORT option.
func (so *SocketOptions) GetReusePort() bool {
	return atomic.LoadUint32(&so.reusePortEnabled) != 0
}

// SetReusePort sets value for SO_REUSEPORT option.
func (so *SocketOptions) SetReusePort(v bool) {
	storeAtomicBool(&so.reusePortEnabled, v)
	so.handler.OnReusePortSet(v)
}

// GetKeepAlive gets value for SO_KEEPALIVE option.
func (so *SocketOptions) GetKeepAlive() bool {
	return atomic.LoadUint32(&so.keepAliveEnabled) != 0
}

// SetKeepAlive sets value for SO_KEEPALIVE option.
func (so *SocketOptions) SetKeepAlive(v bool) {
	storeAtomicBool(&so.keepAliveEnabled, v)
	so.handler.OnKeepAliveSet(v)
}

// GetAcceptConn gets value for SO_ACCEPTCONN option.
func (so *SocketOptions) GetAcceptConn() bool {
	// This option is completely endpoint dependent and unsettable.
	return so.handler.IsListening()
}

// GetLastError gets value for SO_ERROR option.
func (so *SocketOptions) GetLastError() *Error {
	return so.handler.LastError()
}

// GetOutOfBandInline gets value for SO_OOBINLINE option.
func (so *SocketOptions) GetOutOfBandInline() bool {
	return true
}

// SetOutOfBandInline sets value for SO_OOBINLINE option. We currently do not
// support disabling this option.
func (so *SocketOptions) SetOutOfBandInline(v bool) {}

// GetLinger gets value for SO_LINGER option.
func (so *SocketOptions) GetLinger() LingerOption {
	so.mu.Lock()
	linger := so.linger
	so.mu.Unlock()
	return linger
}

// SetLinger sets value for SO_LINGER option.
func (so *SocketOptions) SetLinger(linger LingerOption) {
	so.mu.Lock()
	so.linger = linger
	so.mu.Unlock()
}

// GetBindToDevice gets value for SO_BINDTODEVICE option.
func (so *SocketOptions) GetBindToDevice() NICID {
	so.mu.Lock()
	bindToDevice := so.bindToDevice
	so.mu.Unlock()
	return bindToDevice
}

// SetBindToDevice sets value for SO_BINDTODEVICE option.
func (so *SocketOptions) SetBindToDevice(bindToDevice int32) {
	so.mu.Lock()
	so.bindToDevice = NICID(bindToDevice)
	so.mu.Unlock()
}
