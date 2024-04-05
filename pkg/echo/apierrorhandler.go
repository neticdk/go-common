package echo

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neticdk/go-common/pkg/log"
	"github.com/neticdk/go-common/pkg/types"
)

// APIErrorHandler handles errors occurring in the API classes
//
// Based on [echo.DefaultHTTPErrorHandler] but always logging errors and returning "problem" response with content-type
// "application/problem+json". By default the handler will setup a http status code 500 unless the error is either a
// [types.Problem] or [echo.HTTPError] then the status code from the error will be used.
func APIErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	p, ok := err.(*types.Problem)
	if !ok {
		he, ok := err.(*echo.HTTPError)
		if ok {
			p = &types.Problem{
				Status: &he.Code,
				Detail: fmt.Sprintf("%s", he.Message),
				Err:    he.Internal,
			}
		} else {
			p = &types.Problem{
				Type:   "about:blank",
				Status: types.IntPointer(http.StatusInternalServerError),
				Title:  http.StatusText(http.StatusInternalServerError),
			}
		}
	}

	attrs := []any{}
	attrs = append(attrs, slog.String("uri", c.Request().RequestURI))
	attrs = append(attrs, slog.String("method", c.Request().Method))
	if uw := errors.Unwrap(p); uw != nil {
		attrs = append(attrs, log.Error(uw))
	} else {
		attrs = append(attrs, log.Error(p))
	}
	log.FromContext(c.Request().Context()).ErrorContext(c.Request().Context(), "error occurred processing api request", attrs...)

	status := http.StatusInternalServerError
	if p.Status != nil {
		status = *p.Status
	}
	err = c.JSON(status, p)
	if err != nil {
		log.FromContext(c.Request().Context()).WarnContext(c.Request().Context(), "unable to send error response to client", log.Error(err))
	}
}
