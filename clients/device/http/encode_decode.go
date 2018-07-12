package http

import (
	// stdlib
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	// external
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	// project

	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
)

func encodeUnlockRequest(route *mux.Route) kithttp.EncodeRequestFunc {
	return func(_ context.Context, r *http.Request, request interface{}) error {
		var (
			err error
			req = request.(transport.UnlockRequest)
		)

		if r.URL, err = route.Host(r.URL.Host).URL(
			"event_id", req.EventID.String(),
			"device_id", req.DeviceID.String(),
			"code", req.Code,
		); err != nil {
			return err
		}
		if methods, err := route.GetMethods(); err == nil {
			r.Method = methods[0]
		}

		return nil
	}
}

func decodeUnlockResponse(_ context.Context, response *http.Response) (interface{}, error) {
	var res transport.UnlockResponse
	if response.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		errStr := string(b)
		switch errStr {
		case device.ErrorRequireEventID:
			res.Err = device.ErrRequireEventID
		case device.ErrorRequireDeviceID:
			res.Err = device.ErrRequireDeviceID
		case device.ErrorRequireUnlockCode:
			res.Err = device.ErrRequireUnlockCode
		case device.ErrorRepository:
			res.Err = device.ErrRepository
		case device.ErrorEventNotFound:
			res.Err = device.ErrEventNotFound
		case device.ErrorUnlockNotFound:
			res.Err = device.ErrUnlockNotFound
		default:
			return nil, errors.New(errStr)
		}
		return res, nil
	}
	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
