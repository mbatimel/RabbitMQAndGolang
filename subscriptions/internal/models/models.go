package models

import "go/types"

const (
	MASTER  = true
	REPLICA = false
)

type Err400 struct {
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки`
	// @tg example=`content.api.errors.regressionApi.badRequest`
	ErrorText string `json:"errorText"`
	// @tg desc=`Текст ошибки, при ответе`
	AdditionalErrors struct {
		Errors []struct {
			TrKey string `json:"trKey"`
			// @tg example=`{"1": "value one", "2": "value two"}`
			Params map[string]string `json:"params"`
		} `json:"errors"`
	} `json:"additionalErrors"`
}

type Err403 struct {
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки`
	// @tg example=`content.api.errors.regressionApi.accessDenied`
	ErrorText string `json:"errorText"`
	// @tg desc=`Текст ошибки, при ответе, со статус кодом 403, не указывается`
	AdditionalErrors types.Nil `json:"additionalErrors,omitempty"`
}

type Err405 struct {
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки`
	// @tg example=`content.api.errors.regressionApi.methodNotAllowed`
	ErrorText string `json:"errorText"`
	// @tg desc=`Текст ошибки, при ответе, со статус кодом 403, не указывается`
	AdditionalErrors types.Nil `json:"additionalErrors,omitempty"`
}

type Err500 struct {
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки`
	// @tg example=`content.api.errors.regressionApi.internalError`
	ErrorText string `json:"errorText"`
	// @tg desc=`Текст ошибки, при ответе, со статус кодом 403, не указывается`
	AdditionalErrors types.Nil `json:"additionalErrors,omitempty"`
}

type ActiveSubscriptionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data error `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
