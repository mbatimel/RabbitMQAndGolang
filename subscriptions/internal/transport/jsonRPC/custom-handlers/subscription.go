package customhandlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/config"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/errors"
	subscriptions "github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/interfaces"
	"github.com/rs/zerolog/log"
)

const serviceName = "Regression"

func ActiveSubscription(ctx *fiber.Ctx, svc subscriptions.Subscription, limitId int, price int) error {
	var (
		methodName = "ActiveSubscription"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":  "post",
			"path":    "/activesubscription",
			"limitId": limitId,
			"price":   price,
			"service": serviceName,
			"took":    time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	err = svc.ActiveSubscription(ctx.Context(), limitId, price)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, nil, nil)
	return err
}
