// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package externalapi

import "github.com/gofiber/fiber/v2"

type withRedirect interface {
	RedirectTo() string
}

type cookieType interface {
	Cookie() *fiber.Cookie
}
