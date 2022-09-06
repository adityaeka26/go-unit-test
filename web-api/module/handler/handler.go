package handler

import (
	"net/http"

	"go-unit-test/web-api/helper"
	"go-unit-test/web-api/logger"
	"go-unit-test/web-api/module/model/web"
	"go-unit-test/web-api/module/usecase"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Index(c *gin.Context)
	Register(c *gin.Context)
	VerifyRegister(c *gin.Context)
}
type HandlerImpl struct {
	usecase usecase.Usecase
	logger  logger.Logger
}

func NewHandler(usecase usecase.Usecase, logger logger.Logger) Handler {
	return &HandlerImpl{
		usecase: usecase,
		logger:  logger,
	}
}

func (handler HandlerImpl) Index(c *gin.Context) {
	helper.RespSuccess(c, nil, "Index success")
}

func (handler HandlerImpl) Register(c *gin.Context) {
	request := &web.RegisterRequest{}

	if err := c.ShouldBind(request); err != nil {
		helper.RespError(c, helper.CustomError(http.StatusBadRequest, err.Error()))
		return
	}
	if err := helper.Validate(request); err != nil {
		helper.RespError(c, helper.CustomError(http.StatusBadRequest, err.Error()))
		return
	}

	err := handler.usecase.Register(c.Request.Context(), *request)
	if err != nil {
		helper.RespError(c, err)
		return
	}
	helper.RespSuccess(c, nil, "Register success")
}

func (handler HandlerImpl) VerifyRegister(c *gin.Context) {
	request := &web.VerifyRegisterRequest{}

	if err := c.ShouldBind(request); err != nil {
		helper.RespError(c, helper.CustomError(http.StatusBadRequest, err.Error()))
		return
	}
	if err := helper.Validate(request); err != nil {
		helper.RespError(c, helper.CustomError(http.StatusBadRequest, err.Error()))
		return
	}

	response, err := handler.usecase.VerifyRegister(c.Request.Context(), *request)
	if err != nil {
		helper.RespError(c, err)
		return
	}
	helper.RespSuccess(c, response, "Register success")
}
