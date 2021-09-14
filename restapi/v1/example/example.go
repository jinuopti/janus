package example

import (
    . "github.com/jinuopti/janus/restapi/v1"
    "github.com/labstack/echo/v4"
)

type Request struct {
    Message     string
}

type Response struct {
    Result      string
}

func SetRoute(e *echo.Echo) {
    e.GET(ApiPath + "/example", ExampleHandler)
}

// ExampleHandler
//
// @Summary Example API
// @Description Example API
// @Tags Example
// @Accept json
// @Produce json
// @Security ApiKeyAuth
//
// @Param message query string true "Example message"
//
// @Success 200 {object} JSONResult{data=Response} "Success"
// @Failure 400 {object} JSONResult{data=string} "Error"
// @Failure 401 {object} JSONResult{data=string} "Unauthorized"
// @Router /example [get]
func ExampleHandler(c echo.Context) error {
    mesg := c.FormValue("message")
    if len(mesg) <= 0 {
        return exampleError(c, CodeError, "Error")
    }
    return exampleSuccess(c, "Success")
}

// exampleSuccess API 요청 처리 성공 응답
func exampleSuccess(c echo.Context, message string) error {
    return SuccessResponse(c, Response{Result: message})
}

// exampleError API 요청 처리 실패  응답
func exampleError(c echo.Context, code int, message string) error {
    return ErrorResponse(c, code, Response{Result: message})
}
