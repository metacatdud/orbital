package config

import "errors"

var (
	ErrSecretKeyLength  = errors.New("secret key must be 32bytes long")
	ErrIpRequired       = errors.New("ip required")
	ErrDataPathRequired = errors.New("data path required")
	ErrConfigSave       = errors.New("cannot save config")
	ErrConfigWrite      = errors.New("cannot write config to file")
	ErrConfigRead       = errors.New("cannot read config")
	ErrConfigClient     = errors.New("node cannot be set to client")

	ErrAddrIsEmpty     = errors.New("addr cannot be empty")
	ErrAddrInvalidIP   = errors.New("invalid ip address")
	ErrAddrInvalidIPv6 = errors.New("invalid ipv6 address")
	ErrAddrPort        = errors.New("invalid port")
	ErrAddrPortNaN     = errors.New("port is not a number")
)
