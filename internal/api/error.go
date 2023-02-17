package api

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"net/http"
)

func ResponseFromServiceError(serviceError base.ServiceError, trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(serviceError.Code),
		Blame:      serviceError.Blame,
		TrackingID: trackingID,
		Message:    serviceError.Message,
	}
}

func GeneralParsingError(trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(http.StatusBadRequest),
		Blame:      base.BlameUser,
		TrackingID: trackingID,
		Message:    "failed to parse request parameters",
	}
}

func GeneralSortError(trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(http.StatusBadRequest),
		Blame:      base.BlameUser,
		TrackingID: trackingID,
		Message:    "bad sort parameters",
	}
}

func GeneralPaginationError(trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(http.StatusBadRequest),
		Blame:      base.BlameUser,
		TrackingID: trackingID,
		Message:    "bad pagination parameters",
	}
}

func GeneralFilterError(trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(http.StatusBadRequest),
		Blame:      base.BlameUser,
		TrackingID: trackingID,
		Message:    "bad filter parameters",
	}
}
func GeneralUnexpectedError(trackingID string) base.ResponseFailure {
	return base.ResponseFailure{
		Status:     http.StatusText(http.StatusInternalServerError),
		Blame:      base.BlameUnknown,
		TrackingID: trackingID,
		Message:    "internal error",
	}
}
