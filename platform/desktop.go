//go:build !wasm

package platform

import (
	"github.com/denisbrodbeck/machineid"
)

func GetDeviceId() (string, bool) {
	const appID = "ebiman"
	id, err := machineid.ProtectedID(appID)
	// TODO - log error?
	return id, err == nil
}
