package middlewares

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"im/pkg/errors"
	"im/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

const (
	defaultLimit     = "50"
	defaultPage      = "1"
	defaultSort      = "id:asc"
	defaultSortOrder = "asc"
)

var (
	// ErrInvalidSortOption ...
	ErrInvalidSortOption = errors.New("invalid sort option")
)

// Pagination ...
type Pagination struct {
	Limit     int    `json:"limit" form:"limit" mapstructure:"limit" binding:"omitempty,gte=1,lte=500" minimum:"1" maximum:"500" default:"30"` // 每页数量
	Page      int    `json:"page" form:"page" mapstructure:"page" binding:"omitempty,gte=1" minimum:"1" default:"30"`                          // 页数
	SortBy    string `json:"sortBy" form:"sortBy" mapstructure:"sortBy" binding:"omitempty"`                                                   // 排序字段
	SortOrder string `json:"sortOrder" form:"sortOrder" mapstructure:"sortOrder" binding:"omitempty,oneof=asc desc" enums:"asc,desc"`          // 排序方式
	MultiSort string // 多字段排序
}

func (p *Pagination) SetDefault() {
	p.Limit = 30
	p.Page = 1
	p.SortBy = "id"
	p.SortOrder = "asc"
}

// GetDocsAndTotal ...
func (p *Pagination) GetDocsAndTotal(query *gorm.DB, data interface{}, total *int64) error {
	err := query.
		Count(total).
		Error
	if err != nil {
		return err
	}
	if p.Page > 0 && p.Limit > 0 {
		query = query.
			Offset((p.Page - 1) * p.Limit).
			Limit(p.Limit)
	}
	if len(p.MultiSort) > 0 {
		query = query.Order(p.MultiSort)
	} else {
		if p.SortBy != "" && p.SortOrder != "" {
			if p.SortOrder != "asc" && p.SortOrder != "desc" {
				return ErrInvalidSortOption
			}
			query = query.Order(p.SortBy + " " + p.SortOrder)
		}
	}

	err = query.
		Scan(data).
		Error
	if err != nil {
		return err
	}

	return nil
}

// ScanWhereQuery 检索请求参数拼接where，支持：默认单字段where等于 | 多字段whereOr | 自定义查询方法例如like > <
// example: 表单key：name|number|idCardNo|mobile#like 表单val：123 处理后：
//
//	`name` like '%4442322%' OR `number` like '%4442322%' OR `id_card_no` like '%4442322%' OR `mobile` like '%4442322%'
func (p *Pagination) ScanWhereOrQuery(c *gin.Context, fields ...string) []string {
	where := make([]string, 0)
	fieldVal := ""
	for _, field := range fields {
		if strings.HasSuffix(field, "[]") {
			// 兼容前端框架
			fieldVal = strings.Join(c.QueryArray(field), "|")
			field = strings.TrimRight(field, "[]") + "#in"
		} else {
			fieldVal = c.DefaultQuery(field, "")
			if fieldVal == "" {
				fieldVal = c.DefaultPostForm(field, "")
			}
		}

		if fieldVal != "" {
			fieldName := utils.ToSnakeCase(field)
			sql := ""
			method := "="
			// scan custom method
			customMethod := strings.Split(fieldName, "#")
			if len(customMethod) > 1 {
				fieldName = customMethod[0]
				method = customMethod[1]
			}
			if method == "like" {
				fieldVal = "%" + fieldVal + "%"
			}
			// scan whereOr
			whereOrFields := strings.Split(fieldName, "|")
			for index, thisField := range whereOrFields {
				switch method {
				case "ojs":
					orValues := strings.Split(fieldVal, "|")
					fn := strings.Split(thisField, ":")
					var tql string
					for i, fs := range orValues {
						tql2 := fmt.Sprintf("json_search(`%s`,'one','%s',null,'$[*].%s') is not null", fn[0], fs, strcase.ToLowerCamel(fn[1]))
						if i > 0 {
							tql += " OR " + tql2
						} else {
							tql += tql2
						}
					}
					if index > 0 {
						sql += " OR " + tql
					} else {
						sql += tql
					}

				case "in":
					fv := strings.Split(fieldVal, "|")
					fieldVal = "('" + strings.Join(fv, "','") + "')"
					sql = fmt.Sprintf("`%s` %s %s", thisField, method, fieldVal)
				default:
					fieldVal = strings.ReplaceAll(fieldVal, "'", "")
					tempSql := fmt.Sprintf("`%s` %s '%s'", thisField, method, fieldVal)
					if index > 0 {
						tempSql = " OR " + tempSql
					}
					sql += tempSql
				}
			}
			where = append(where, sql)
		}
	}
	return where
}

func (p *Pagination) ScanWhereDateRange(c *gin.Context, field string, toTimestamp bool) []string {
	where := make([]string, 0)
	dateRange := c.QueryArray("dateRange[]")
	if len(dateRange) != 2 {
		return where
	}
	start := dateRange[0]
	end := dateRange[1]
	startArr := strings.Split(start, " ")
	endArr := strings.Split(end, " ")
	if len(startArr) < 1 || len(endArr) < 1 {
		return where
	}

	startTime := startArr[0] + " 00:00:00"
	endTime := endArr[0] + " 23:59:59"
	if toTimestamp {
		where = append(where, fmt.Sprintf("`%s` >= %d", utils.ToSnakeCase(field), utils.StringToTimestamp(startTime)))
		where = append(where, fmt.Sprintf("`%s` <= %d", utils.ToSnakeCase(field), utils.StringToTimestamp(endTime)))
	} else {
		where = append(where, fmt.Sprintf("`%s` >= '%s'", utils.ToSnakeCase(field), startTime))
		where = append(where, fmt.Sprintf("`%s` <= '%s'", utils.ToSnakeCase(field), endTime))
	}
	return where
}

func (p *Pagination) InsertQueryWhere(query *gorm.DB, where []string) *gorm.DB {
	if len(where) > 0 {
		for _, v := range where {
			query = query.Where(v)
		}
	}
	return query
}

// Response ...
func (p *Pagination) Response(c *gin.Context, total uint64, data interface{}) {

	c.JSON(http.StatusOK, gin.H{
		"items": data,
		"total": total,
	})
}

// NewPaginationMiddleware ...
func NewPaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := &Pagination{}

		limitQuery := c.DefaultQuery("limit", defaultLimit)
		pageQuery := c.DefaultQuery("page", defaultPage)

		ms := strings.Split(c.DefaultQuery("sort", defaultSort), "|")
		if len(ms) > 1 && len(ms[0]) > 0 {
			var mo []string
			for _, m := range ms {
				sorts := strings.Split(m, ":")
				so := defaultSortOrder
				if len(sorts) == 1 || (len(sorts) > 1 && sorts[1] != defaultSortOrder) {
					so = "desc"
				}
				sk := utils.ToSnakeCase(sorts[0])
				sk = "`" + strings.ReplaceAll(sk, "`", "") + "`"
				mo = append(mo, sk+" "+so)
			}

			p.MultiSort = strings.Join(mo, ",")
		} else if len(ms[0]) > 0 {
			sorts := strings.Split(ms[0], ":")
			p.SortOrder = defaultSortOrder
			if len(sorts) > 0 && sorts[1] != defaultSortOrder {
				p.SortOrder = "desc"
			}
			p.SortBy = utils.ToSnakeCase(sorts[0])
			p.SortBy = "`" + strings.ReplaceAll(p.SortBy, "`", "") + "`"
		} else {
			_ = c.Error(ErrInvalidSortOption).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		limit, err := strconv.Atoi(limitQuery)
		if err != nil {
			_ = c.Error(err).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		p.Limit = int(math.Max(1, math.Min(500, float64(limit))))

		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			_ = c.Error(err).SetType(gin.ErrorTypePublic)
			c.Abort()
			return
		}

		p.Page = int(math.Max(1, float64(page)))

		c.Set("pagination", p)
		c.Next()
	}
}
