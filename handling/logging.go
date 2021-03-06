package handling

import (
	"time"

	"github.com/go-kit/kit/log"

	"github.com/marcusolsson/goddd/cargo"
	"github.com/marcusolsson/goddd/location"
	"github.com/marcusolsson/goddd/voyage"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService ...
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) RegisterHandlingEvent(completionTime time.Time, trackingID cargo.TrackingID, voyageNumber voyage.Number,
	unLocode location.UNLocode, eventType cargo.HandlingEventType) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "register_incident",
			"tracking_id", trackingID,
			"location", unLocode,
			"voyage", voyageNumber,
			"event_type", eventType,
			"completion_time", completionTime,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.RegisterHandlingEvent(completionTime, trackingID, voyageNumber, unLocode, eventType)
}
