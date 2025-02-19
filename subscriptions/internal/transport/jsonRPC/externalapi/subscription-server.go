// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package externalapi

import (
	"context"
	"github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/interfaces"
)

type serverSubscription struct {
	svc                interfaces.Subscription
	activeSubscription SubscriptionActiveSubscription
}

type MiddlewareSetSubscription interface {
	Wrap(m MiddlewareSubscription)
	WrapActiveSubscription(m MiddlewareSubscriptionActiveSubscription)

	WithMetrics()
	WithLog()
}

func newServerSubscription(svc interfaces.Subscription) *serverSubscription {
	return &serverSubscription{
		activeSubscription: svc.ActiveSubscription,
		svc:                svc,
	}
}

func (srv *serverSubscription) Wrap(m MiddlewareSubscription) {
	srv.svc = m(srv.svc)
	srv.activeSubscription = srv.svc.ActiveSubscription
}

func (srv *serverSubscription) ActiveSubscription(ctx context.Context, limitId int, price int) (err error) {
	return srv.activeSubscription(ctx, limitId, price)
}

func (srv *serverSubscription) WrapActiveSubscription(m MiddlewareSubscriptionActiveSubscription) {
	srv.activeSubscription = m(srv.activeSubscription)
}

func (srv *serverSubscription) WithMetrics() {
	srv.Wrap(metricsMiddlewareSubscription)
}

func (srv *serverSubscription) WithLog() {
	srv.Wrap(loggerMiddlewareSubscription())
}
