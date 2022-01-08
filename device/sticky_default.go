//go:build !linux

package device

import (
	"github.com/bugfan/wireguard-go/conn"
	"github.com/bugfan/wireguard-go/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
