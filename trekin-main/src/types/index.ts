// Market Types
export interface MonthlyTrend {
  date: string;
  revenue?: number;
  volume?: number;
}

export interface MarketMetrics {
  totalRevenue: number;
  totalProducts: number;
  brandCount: number;
  searchVolume: number;
  cagr: number;
  monthlyTrends: MonthlyTrend[];
  keywordMonthlyTrends: MonthlyTrend[];
}

export interface MarketSummary {
  id: string;
  name: string;
  metrics: MarketMetrics;
}

// Product Types
export interface Product {
  asin: string;
  title: string;
  brand: string;
  price: number;
  rating: number;
  reviews: number;
  monthlySales: number;
  imageUrl: string;
}

export interface ProductSales {
  asin: string;
  monthlySales: {
    date: string;
    sales: number;
    units: number;
  }[];
}

// Keyword Types
export interface Keyword {
  keyword: string;
  source: 'google' | 'amazon';
  monthlyVolume: {
    date: string;
    volume: number;
  }[];
}

// Brand Types
export interface SocialMedia {
  url: string;
  subscribers?: number;
  followers?: number;
  posts?: number;
}

export interface BrandMetrics {
  brand: string;
  totalRevenue: number;
  productCount: number;
  cagr: number;
  monthlyTrends: MonthlyTrend[];
  website?: string;
  youtube?: SocialMedia;
  instagram?: SocialMedia;
  facebook?: SocialMedia;
  reddit?: SocialMedia;
}

// Market Detail
export interface MarketDetail {
  id: string;
  name: string;
  metrics: MarketMetrics;
  products: Product[];
  productSales: { [asin: string]: ProductSales };
  keywords: Keyword[];
  brands: BrandMetrics[];
}

