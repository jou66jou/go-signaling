package pagination

const globalDefaultPerPage = 20

// PageMeta meta with pagination
type PageMeta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination 用來表示分頁
type Pagination struct {
	Page       int32 `query:"page" form:"page" json:"page" description:"目前頁面"`
	PerPage    int32 `query:"perPage" form:"perPage" json:"perPage" description:"每頁顯示多少筆"`
	TotalCount int32 `query:"totalCount" form:"totalCount" json:"totalCount" description:"總筆數"`
	TotalPage  int32 `query:"totalPage" form:"totalPage" json:"totalPage" description:"總頁數"`
}

// Count 用來儲存sql count 結果
type Count struct {
	Count int64 `json:"count" description:"sql count result"`
}

// SetTotalCountAndPage 用來計算總數和分頁
func (p *Pagination) SetTotalCountAndPage(total int32) {
	p.CheckOrSetDefault()
	p.TotalCount = total

	quotient := p.TotalCount / p.PerPage
	remainder := p.TotalCount % p.PerPage
	if remainder > 0 {
		quotient++
	}
	p.TotalPage = quotient
}

// CheckOrSetDefault 檢查Page值若未設置則設置預設值
func (p *Pagination) CheckOrSetDefault(params ...int32) *Pagination {
	var defaultPerPage int32
	if len(params) >= 1 {
		defaultPerPage = params[0]
	}

	if defaultPerPage <= 0 {
		defaultPerPage = globalDefaultPerPage
	}

	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = defaultPerPage
	}
	return p
}

// LimitAndOffset return limit and offset
func (p *Pagination) LimitAndOffset() (uint32, uint32) {
	return uint32(p.PerPage), uint32(p.Offset())
}

// Offset 計算 offset 的值
func (p *Pagination) Offset() int32 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

// APIMetaPagination 分頁資訊
type APIMetaPagination struct {
	Pagination Pagination `json:"pagination" description:"分頁資訊"`
}

// APIDataResult 單存只顯示資料沒有 metadata
type APIDataResult struct {
	Data interface{} `json:"data"`
}

// APIPaginationResult 用來顯示 meta 和 data
type APIPaginationResult struct {
	Meta APIMetaPagination `json:"meta"`
	Data interface{}       `json:"data"`
}
