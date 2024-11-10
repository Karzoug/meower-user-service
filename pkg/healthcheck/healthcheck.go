// Package healthcheck provides a standart types and method for
// healthcheck http server and its dependencies or sub-component,
// inspired mostly by the upcoming IETF RFC Health Check Response Format for HTTP APIs
// https://inadarei.github.io/rfc-healthcheck/ .
package healthcheck

import (
	"context"
	"errors"

	"github.com/Karzoug/meower-user-service/pkg/ucerr"
)

type HealthChecker interface {
	Healthcheck(ctx context.Context) (key string, err error)
}

type status string

const (
	// Pass represents a healthy service.
	Pass status = "pass"
	// fail represents an unhealthy service.
	Fail status = "fail"
	// Warn represents a healthy service with some minor problem.
	Warn status = "warn"
)

type ComponentCheckResponse struct {
	Status status `json:"status"`
	Output string `json:"output,omitempty"`
}

type Response struct {
	Status status                            `json:"status"`
	Checks map[string]ComponentCheckResponse `json:"checks"`
}

type healthcheck struct {
	checkers []HealthChecker
}

func New(checkers ...HealthChecker) *healthcheck {
	return &healthcheck{
		checkers: checkers,
	}
}

func (hc *healthcheck) HealthCheck(ctx context.Context) Response {
	resp := Response{
		Status: Pass,
		Checks: make(map[string]ComponentCheckResponse),
	}

	var hasProblems bool
	for _, checker := range hc.checkers {
		key, err := checker.Healthcheck(ctx)
		cr := ComponentCheckResponse{
			Status: Pass,
		}
		if err != nil {
			var serr ucerr.Error
			if errors.As(err, &serr) {
				cr.Status = Fail
				cr.Output = serr.Error()
			} else {
				cr.Status = Fail
			}
			hasProblems = true
		}
		resp.Checks[key] = cr
	}

	if hasProblems {
		resp.Status = Fail
	}

	return resp
}
