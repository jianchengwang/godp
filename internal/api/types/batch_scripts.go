package types

import "godp/pkg/api/page"

type BatchScriptsQueryParam struct {
	page.PageInfo
	Q      string `form:"q"`
	UserId uint   `form:"userId"`
}
