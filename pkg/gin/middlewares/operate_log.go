package middlewares

import (
	"im/models"
	"im/pkg/database"
	"im/pkg/utils"
	"im/services/system/operate_log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewOperateLogger(log *zap.SugaredLogger, mysqlClient *database.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		go func() {
			op, exists := c.Get(operate_log.KEY_LOG)
			if !exists {
				return
			}
			operateLog, ok := op.(*models.OperateLogs)
			if !ok {
				log.Errorf("op.(*models.OperateLogs) err")
				return
			}
			operateLog.IP = utils.ConvIP(c.Request.Header)
			operateLog.AccountID = c.GetUint("id")

			needFields := []string{}
			if val, exists := c.Get(operate_log.KEY_NEED_FIELDS); exists {
				needFields, ok = val.([]string)
				if !ok {
					log.Errorf("needFields get err")
					return
				}
			}
			expectFields := []string{}
			if val, exists := c.Get(operate_log.KEY_EXPECT_FIELDS); exists {
				expectFields, ok = val.([]string)
				if !ok {
					log.Errorf("existsFields get err")
					return
				}
			}

			beforeFields, existsBefore := c.Get(operate_log.KEY_BEFORE)
			afterFields, existsAfter := c.Get(operate_log.KEY_AFTER)

			if existsBefore || existsAfter {
				fieldsSlice, beforeSlice, afterSlice, err := operate_log.GetFieldsLogSlice(beforeFields, afterFields, needFields, expectFields)
				if err != nil {
					log.Errorf("GetFieldsLogSlice err %v", err)
					return
				}
				operateLog.Fields = fieldsSlice
				operateLog.FieldsBefore = beforeSlice
				operateLog.FieldsAfter = afterSlice
			}

			err := mysqlClient.Db().
				Create(operateLog).
				Error
			if err != nil {
				log.Errorf("create operateLog %v", err)
				return
			}
		}()
	}
}
