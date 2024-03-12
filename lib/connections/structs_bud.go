package connections

import (
	"errors"
	"fmt"
)

const (
	// connTypeRelayClient is a connType of type relay-client.
	connTypeRelayClient connType = iota
	// connTypeRelayServer is a connType of type relay-server.
	connTypeRelayServer
	// connTypeTCPClient is a connType of type tcp-client.
	connTypeTCPClient
	// connTypeTCPServer is a connType of type tcp-server.
	connTypeTCPServer
	// connTypeQUICClient is a connType of type quic-client.
	connTypeQUICClient
	// connTypeQUICServer is a connType of type quic-server.
	connTypeQUICServer
)

var ErrInvalidconnType = errors.New("not a valid connType")

var _connTypeName = "relay-clientrelay-servertcp-clienttcp-serverquic-clientquic-server"

var _connTypeMapName = map[connType]string{
	connTypeRelayClient: _connTypeName[0:12],
	connTypeRelayServer: _connTypeName[12:24],
	connTypeTCPClient:   _connTypeName[24:34],
	connTypeTCPServer:   _connTypeName[34:44],
	connTypeQUICClient:  _connTypeName[44:55],
	connTypeQUICServer:  _connTypeName[55:66],
}

// Name is the attribute of connType.
func (x connType) Name() string {
	if v, ok := _connTypeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("connType(%d).Name", x)
}

var _connTypeMapTransport = map[connType]string{
	connTypeRelayClient: "relay",
	connTypeRelayServer: "relay",
	connTypeTCPClient:   "tcp",
	connTypeTCPServer:   "tcp",
	connTypeQUICClient:  "quic",
	connTypeQUICServer:  "quic",
}

// Transport is the attribute of connType.
func (x connType) Transport() string {
	if v, ok := _connTypeMapTransport[x]; ok {
		return v
	}
	return fmt.Sprintf("connType(%d).Transport", x)
}

// Val is the attribute of connType.
func (x connType) Val() int {
	return int(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x connType) IsValid() bool {
	_, ok := _connTypeMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x connType) String() string {
	return x.Name()
}

var _connTypeNameMap = map[string]connType{
	_connTypeName[0:12]:  connTypeRelayClient,
	_connTypeName[12:24]: connTypeRelayServer,
	_connTypeName[24:34]: connTypeTCPClient,
	_connTypeName[34:44]: connTypeTCPServer,
	_connTypeName[44:55]: connTypeQUICClient,
	_connTypeName[55:66]: connTypeQUICServer,
}

// ParseconnType converts a string to a connType.
func ParseconnType(value string) (connType, error) {
	if x, ok := _connTypeNameMap[value]; ok {
		return x, nil
	}
	return connType(0), fmt.Errorf("%s is %w", value, ErrInvalidconnType)
}
