// Package interfaces
// @tg version=0.0.1
// @tg backend=subscription
// @tg title=`Subscriptions API`
//
//go:generate tg transport --services . --out ../../internal/transport/jsonRPC/externalapi --outSwagger ../../swaggers/subscription/swagger.yaml
package interfaces

import (
	"context"
)

// subscription
// @tg http-server metrics log
// @tg http-prefix=/api/v1
// @tg 200=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:Resp200
// @tg 403=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:Err403
// @tg 405=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:Err405
// @tg 500=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:Err500
type Subscription interface {
	// ActiveSubscription ...
	// @tg http-method=POST
	// @tg http-path=/activesubscription
	// @tg summary=`Ручка активации подписки`
	// @tg http-response=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/transport/jsonRPC/custom-handlers:ActiveSubscription
	// @tg desc=`Ручка возвращает bool результут подкиски`
	// @tg 400=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:Err400
	// @tg 200=github.com/mbatimel/RabbitMQAndGolang/subscriptions/internal/models:ActiveSubscriptionResp200
	ActiveSubscription(ctx context.Context, limitId int, price int) (err error)
}
