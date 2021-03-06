package functions

import (
	"context"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// Call simply calls a function on a service
func (m *Module) Call(service, function, token string, reqParams model.RequestParams, params interface{}, timeout int) (int, interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	status, result, err := m.handleCall(ctx, service, function, token, reqParams.Claims, params)
	if err != nil {
		return status, result, err
	}
	m.metricHook(m.project, service, function)
	return status, result, nil
}

// CallWithContext invokes function on a service. The response from the function is returned back along with
// any errors if they occurred.
func (m *Module) CallWithContext(ctx context.Context, service, function, token string, reqParams model.RequestParams, params interface{}) (int, interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	status, result, err := m.handleCall(ctx, service, function, token, reqParams.Claims, params)
	if err != nil {
		return status, result, err
	}
	m.metricHook(m.project, service, function)
	return status, result, nil
}

// AddInternalRule add an internal rule to internal service
func (m *Module) AddInternalRule(service string, rule *config.Service) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.config.InternalServices[service] = rule
}
