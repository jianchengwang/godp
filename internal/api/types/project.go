package types

import "godp/pkg/api/page"

type ProjectQueryParam struct {
	page.PageInfo
	Q        string `form:"q"`
	Internal string `form:"internal"`
}
