package middlewares

import (
	"net/http"

	"im/pkg/errors"
	"im/pkg/resp"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

var (
	matcher = language.NewMatcher([]language.Tag{
		language.English,
		language.AmericanEnglish,
		language.BritishEnglish,
		language.Chinese,
		language.SimplifiedChinese,
		language.TraditionalChinese,
	})
)

// ErrorHandler ...
type ErrorHandler struct {
	uni *ut.UniversalTranslator
	log *zap.SugaredLogger
}

// NewWithStatusHandler ...
func NewWithStatusHandler(uni *ut.UniversalTranslator) *ErrorHandler {
	return &ErrorHandler{
		uni: uni,
		log: zap.S().With("module", "errorHandlers"),
	}
}

// HandleErrors ...
func (h *ErrorHandler) HandleErrors(c *gin.Context) {
	c.Next()

	errorToPrint := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if errorToPrint == nil {
		return
	}
	h.log.Errorf("Error captured, stacktrace: %+v", errorToPrint.Err)

	al := c.Request.Header.Get("Accept-Language")
	lang, _ := c.Request.Cookie("lang")

	tag, _ := language.MatchStrings(matcher, lang.String(), al)
	trans, _ := h.uni.GetTranslator(tag.String())

	if _, ok := errorToPrint.Err.(validator.ValidationErrors); ok {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"message": "validation error",
		// 	"errors":  errs.Translate(trans),
		// })
		c.JSON(http.StatusOK, resp.Response{
			Code:    1,
			Message: resp.PARAM_INVALID,
			Result:  errorToPrint.Error(),
		})
		return
	}

	if getter, ok := errorToPrint.Err.(errors.StatusCodeGetter); ok {
		msg := errorToPrint.Error()
		if t, ok := errorToPrint.Err.(errors.Translator); ok {
			msg = t.Translate(trans)
		}
		data := resp.Response{
			Code:    1,
			Message: msg,
			Result:  msg,
		}
		// if t, ok := errorToPrint.Err.(errors.DetailGetter); ok {
		// 	data["details"] = t.Details()
		// }
		c.JSON(getter.HTTPStatusCode(), data)
		return
	}

	// status := http.StatusInternalServerError
	// if errorToPrint.Meta != nil {
	// 	status = errorToPrint.Meta.(gin.H)["status"].(int)
	// 	if status == 0 {
	// 		status = http.StatusInternalServerError
	// 	}
	// }
	c.JSON(http.StatusOK, resp.Response{
		Code:    1,
		Message: errorToPrint.Error(),
		Result:  errorToPrint.Error(),
	})
}
