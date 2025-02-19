// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package externalapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/interfaces"
)

type httpSubscription struct {
	errorHandler     ErrorHandler
	maxBatchSize     int
	maxParallelBatch int
	svc              *serverSubscription
	base             interfaces.Subscription
}

func NewSubscription(svcSubscription interfaces.Subscription) (srv *httpSubscription) {

	srv = &httpSubscription{
		base: svcSubscription,
		svc:  newServerSubscription(svcSubscription),
	}
	return
}

func (http *httpSubscription) Service() *serverSubscription {
	return http.svc
}

func (http *httpSubscription) WithLog() *httpSubscription {
	http.svc.WithLog()
	return http
}

func (http *httpSubscription) WithMetrics() *httpSubscription {
	http.svc.WithMetrics()
	return http
}

func (http *httpSubscription) WithErrorHandler(handler ErrorHandler) *httpSubscription {
	http.errorHandler = handler
	return http
}

func (http *httpSubscription) SetRoutes(route *fiber.App) {
	route.Post("/api/v1/activesubscription", http.serveActiveSubscription)
}
