// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package externalapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

const logLevelHeader = "X-Log-Level"

func (srv *Server) setLogger(ctx *fiber.Ctx) error {
	ctx.SetUserContext(srv.log.WithContext(ctx.UserContext()))
	return ctx.Next()
}

func (srv *Server) logLevelHandler(ctx *fiber.Ctx) error {

	if levelName := string(ctx.Request().Header.Peek(logLevelHeader)); levelName != "" {
		if level, err := zerolog.ParseLevel(levelName); err == nil {
			logger := log.Ctx(ctx.UserContext()).Level(level)
			ctx.SetUserContext(logger.WithContext(ctx.UserContext()))
		}
	}
	return ctx.Next()
}

func recoverHandler(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.New(fmt.Sprintf("%v", r))
			}
			log.Ctx(ctx.UserContext()).Error().Stack().Err(errors.Wrap(err, "recover")).Str("method", ctx.Method()).Str("path", ctx.OriginalURL()).Msg("panic occurred")
			ctx.Status(fiber.StatusInternalServerError)
		}
	}()
	return ctx.Next()
}
