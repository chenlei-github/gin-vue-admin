package response

//import "time"

// MonthlyTrend 月度趋势数据
type MonthlyTrend struct {
	Date         string  `json:"date"` // YYYY-MM格式
	Revenue      float64 `json:"revenue"`
	SearchVolume *int    `json:"searchVolume,omitempty"` // 仅市场级别有
}

// ProductSalesTrend 商品销售趋势
type ProductSalesTrend struct {
	Date  string  `json:"date"` // YYYY-MM格式
	Sales float64 `json:"sales"`
}

// SocialMedia 社交媒体信息
type SocialMedia struct {
	Youtube   *SocialPlatform `json:"youtube,omitempty"`
	Instagram *SocialPlatform `json:"instagram,omitempty"`
	Facebook  *SocialPlatform `json:"facebook,omitempty"`
	Reddit    *SocialPlatform `json:"reddit,omitempty"`
}

// SocialPlatform 社交平台
type SocialPlatform struct {
	URL         string `json:"url"`
	Subscribers *int   `json:"subscribers,omitempty"` // YouTube
	Followers   *int   `json:"followers,omitempty"`   // Instagram/Facebook
	Posts       *int   `json:"posts,omitempty"`       // Reddit
}

// MarketMetrics 市场指标
type MarketMetrics struct {
	TotalRevenue  float64        `json:"totalRevenue"`
	TotalProducts int            `json:"totalProducts"`
	BrandCount    int            `json:"brandCount"`
	SearchVolume  int            `json:"searchVolume"`
	CAGR          *float64       `json:"cagr"`
	MonthlyTrends []MonthlyTrend `json:"monthlyTrends"`
}

// MarketListItem 市场列表项
type MarketListItem struct {
	ID      string        `json:"id"`   // market_slug
	Name    string        `json:"name"` // market_name
	Metrics MarketMetrics `json:"metrics"`
}

// BrandSummary 品牌摘要（用于市场详情中的品牌列表）
type BrandSummary struct {
	Brand        string       `json:"brand"`
	TotalRevenue float64      `json:"totalRevenue"`
	ProductCount int          `json:"productCount"`
	CAGR         *float64     `json:"cagr"`
	Website      *string      `json:"website,omitempty"`
	Social       *SocialMedia `json:"social,omitempty"`
}

// MarketDetail 市场详情
type MarketDetail struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Metrics MarketMetrics  `json:"metrics"`
	Brands  []BrandSummary `json:"brands"`
}

// ProductDetail 商品详情
type ProductDetail struct {
	ASIN         string              `json:"asin"`
	Title        string              `json:"title"`
	Price        float64             `json:"price"`
	Rating       float64             `json:"rating"`
	Reviews      int                 `json:"reviews"`
	ImageUrl     string              `json:"imageUrl"`
	MonthlySales int                 `json:"monthlySales"`
	SalesTrend   []ProductSalesTrend `json:"salesTrend"`
}

// BrandDetail 品牌详情
type BrandDetail struct {
	Brand         string          `json:"brand"`
	TotalRevenue  float64         `json:"totalRevenue"`
	ProductCount  int             `json:"productCount"`
	CAGR          *float64        `json:"cagr"`
	Website       *string         `json:"website,omitempty"`
	Social        *SocialMedia    `json:"social,omitempty"`
	MonthlyTrends []MonthlyTrend  `json:"monthlyTrends"`
	Products      []ProductDetail `json:"products"`
}
