package page

type PageInfo struct {
	CurPage int                    `form:"curPage"`
	Limit   int                    `form:"limit"`
	Sidx    string                 `form:"sidx"`
	Order   string                 `form:"order"`
	Filter  map[string]interface{} `form:"filter"`
}

type PageResult struct {
	PageSize   int         `json:"pageSize"`
	TotalPage  int         `json:"totalPage"`
	CurrPage   int         `json:"currPage"`
	TotalCount int64       `json:"totalCount"`
	List       interface{} `json:"list"`
}

func PageResultWrapper(curPage int, pageSize int, objList interface{}, total int64) PageResult {
	var pageResult = PageResult{
		PageSize:   pageSize,
		TotalPage:  int(total) / pageSize,
		CurrPage:   curPage,
		TotalCount: total,
		List:       objList,
	}
	return pageResult
}
