package service

import (
	"encoding/json"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/danielpaulus/go-ios/ios"
	"strings"
)
func XCTestHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bundleID := strings.TrimSpace(r.URL.Query().Get("bundleid"))
		udid := strings.TrimSpace(r.URL.Query().Get("udid"))
		device, err := ios.GetDevice(udid)
		if err != nil {
			serverError(err.Error(), http.StatusInternalServerError, w)
			return
		}
		err = testmanagerd.RunXCUITest(bundleID, device)
		if err != nil {
			serverError(err.Error(), http.StatusInternalServerError, w)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}}


//HealthHandler is a simple health check. It executes a basic codesign operation to make sure
//codesign really works and is set up correctly
func HealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err:= ios.ListDevices()
		if err != nil {
			serverError("failed getting devicelist", http.StatusInternalServerError, w)
			return
		}
		json, err := json.Marshal(
			map[string]string{
				"version":           GetVersion(),
				"devices": devices.String(),
			},
		)
		if err != nil {
			serverError("failed encoding json", http.StatusInternalServerError, w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)

	}
}





func serverError(message string, code int, w http.ResponseWriter) {
	json, err := json.Marshal(
		map[string]string{"error": message},
	)
	if err != nil {
		log.Warnf("error encoding json:%+v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}
