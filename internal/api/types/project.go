package types

import "godp/pkg/page"

type ProjectQueryParam struct {
	page.PageInfo
	Q        string `form:"q"`
	Internal string `form:"internal"`
}
