package web

import (
	"fmt"

	"github.com/ca17/teamsacs/common"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PreQuery struct {
	context          echo.Context
	defaultOrderby   string
	dateRange        DateRange
	timeField        string
	equalFieldds     []string
	keyfilterFieldds []string
	params           map[string]interface{}
}

func NewPreQuery(c echo.Context) *PreQuery {
	return &PreQuery{context: c, params: make(map[string]interface{})}
}

func (p *PreQuery) DefaultOrderBy(fd string) *PreQuery {
	p.defaultOrderby = fd
	return p
}

func (p *PreQuery) DateRange(queryfd, timefd string) *PreQuery {
	daterange, err := NewWebForm(p.context).GetDateRange(queryfd)
	if err == nil {
		p.dateRange = daterange
		p.timeField = timefd
	}
	return p
}

func (p *PreQuery) EqualFields(fd ...string) *PreQuery {
	p.equalFieldds = fd
	return p
}

func (p *PreQuery) KeyFields(fd ...string) *PreQuery {
	p.keyfilterFieldds = fd
	return p
}

func (p *PreQuery) SetParam(key string, value interface{}) *PreQuery {
	p.params[key] = value
	return p
}

func (p *PreQuery) Query(query *gorm.DB) *gorm.DB {
	if len(ParseSortMap(p.context)) == 0 {
		query = query.Order(p.defaultOrderby)
	} else {
		for name, stype := range ParseSortMap(p.context) {
			query = query.Order(fmt.Sprintf("%s %s", name, stype))
		}
	}

	if _start, err := p.dateRange.ParseStart(); err == nil {
		query = query.Where(p.timeField+" >= ? ", _start)
	}
	if _end, err := p.dateRange.ParseEnd(); err == nil {
		query = query.Where(p.timeField+" <= ?", _end)
	}

	for name, value := range ParseEqualMap(p.context) {
		if common.IsEmptyOrNA(value) {
			continue
		}
		if _, ok := p.params[name]; ok {
			continue
		}
		query = query.Where(fmt.Sprintf("%s = ?", name), value)
	}

	for name, value := range p.params {
		query = query.Where(fmt.Sprintf("%s = ?", name), value)
	}

	for name, value := range ParseFilterMap(p.context) {
		if common.IsEmptyOrNA(value) {
			continue
		}
		if _, ok := p.params[name]; ok {
			continue
		}
		if common.InSlice(name, p.equalFieldds) {
			query = query.Where(fmt.Sprintf("%s = ?", name), value)
		} else {
			query = query.Where(fmt.Sprintf("%s like ?", name), "%"+value+"%")
		}
	}

	keyword := p.context.QueryParam("keyword")
	if keyword != "" {
		for i, keyfd := range p.keyfilterFieldds {
			if i == 0 {
				query = query.Where(keyfd+" like ?", "%"+keyword+"%")
			} else {
				query = query.Or(keyfd+" like ?", "%"+keyword+"%")
			}
		}
	}

	return query
}
