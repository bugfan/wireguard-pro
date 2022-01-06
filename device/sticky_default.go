//go:build !linux

package device

import (
	"github.com/bugfan/wire/conn"
	"github.com/bugfan/wire/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
