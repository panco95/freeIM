package middlewares

import (
	"im/resource"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func loadI18nFiles(bundle *i18n.Bundle) error {
	filenames, err := resource.ReadFilenames("root/i18n")
	if err != nil {
		return err
	}

	for _, filename := range filenames {
		buf, err := resource.ReadAll(filename)
		if err != nil {
			return err
		}

		_, err = bundle.ParseMessageFileBytes(buf, filename)
		if err != nil {
			return err
		}
		zap.S().Infof("Load i18n file %s done", filename)
	}
	return nil
}

// NewI18nMiddleware ...
func NewI18nMiddleware() gin.HandlerFunc {
	bundle := i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	err := loadI18nFiles(bundle)
	if err != nil {
		zap.S().Errorf("Load i18n file error %v", err)
	}

	return func(c *gin.Context) {
		lang, _ := c.Cookie("lang")
		accept := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)
		c.Set("localizer", localizer)
		c.Next()
	}
}
