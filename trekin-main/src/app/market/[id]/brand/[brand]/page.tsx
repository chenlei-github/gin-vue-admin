"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import {
  TrendingUp,
  TrendingDown,
  Package,
  ArrowLeft,
  ChevronRight,
  Youtube,
  Instagram,
  Facebook,
  Globe,
  MessageCircle,
  Star,
  ExternalLink,
} from "lucide-react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { MarketDetail, BrandMetrics, Product } from "@/types";

export default function BrandDetailPage() {
  const params = useParams();
  const [market, setMarket] = useState<MarketDetail | null>(null);
  const [brand, setBrand] = useState<BrandMetrics | null>(null);
  const [brandProducts, setBrandProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (params.id && params.brand) {
      const brandName = decodeURIComponent(params.brand as string);

      fetch(`/data/${params.id}.json`)
        .then((res) => res.json())
        .then((data: MarketDetail) => {
          setMarket(data);

          const brandData = data.brands.find((b) => b.brand === brandName);
          setBrand(brandData || null);

          const products = data.products.filter((p) => p.brand === brandName);
          setBrandProducts(products);

          setLoading(false);
        })
        .catch((err) => {
          console.error("Error loading brand:", err);
          setLoading(false);
        });
    }
  }, [params.id, params.brand]);

  const formatNumber = (num: number) => {
    if (num >= 1e9) return `$${(num / 1e9).toFixed(2)}B`;
    if (num >= 1e6) return `$${(num / 1e6).toFixed(2)}M`;
    if (num >= 1e3) return `$${(num / 1e3).toFixed(2)}K`;
    return `$${num.toFixed(0)}`;
  };

  const formatVolume = (num: number) => {
    if (num >= 1e6) return `${(num / 1e6).toFixed(2)}M`;
    if (num >= 1e3) return `${(num / 1e3).toFixed(1)}K`;
    return num.toFixed(0);
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (!market || !brand) {
    return (
      <div className="container" style={{ padding: "48px 24px", textAlign: "center" }}>
        <h1>品牌未找到</h1>
        <Link href={`/market/${params.id}`}>
          <button className="button button-primary" style={{ marginTop: "24px" }}>
            返回市场页
          </button>
        </Link>
      </div>
    );
  }

  // Prepare trend data - show all available data
  const trendData = brand.monthlyTrends;

  return (
    <div className="container" style={{ padding: "48px 24px" }}>
      {/* Breadcrumb */}
      <div className="breadcrumb">
        <Link href="/">市场分析</Link>
        <ChevronRight size={16} />
        <Link href={`/market/${params.id}`}>{market.name}</Link>
        <ChevronRight size={16} />
        <span>{brand.brand}</span>
      </div>

      {/* Header */}
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "start",
          marginBottom: "32px",
        }}
      >
        <div>
          <h1 className="page-title">{brand.brand}</h1>
          <p className="page-subtitle">
            品牌销售趋势、社交媒体数据和产品信息
          </p>
        </div>
        <Link href={`/market/${params.id}`}>
          <button className="button button-secondary">
            <ArrowLeft size={16} />
            返回市场
          </button>
        </Link>
      </div>

      {/* Key Metrics */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">品牌规模</span>
            <div className="stat-icon" style={{ background: "#dcfce7", color: "#10b981" }}>
              <TrendingUp size={20} />
            </div>
          </div>
          <div className="stat-value">{formatNumber(brand.totalRevenue)}</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            年度销售额预估
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">品牌增速</span>
            <div
              className="stat-icon"
              style={{
                background: brand.cagr !== null && brand.cagr >= 0 ? "#dcfce7" : brand.cagr !== null ? "#fee2e2" : "#f1f5f9",
                color: brand.cagr !== null && brand.cagr >= 0 ? "#10b981" : brand.cagr !== null ? "#ef4444" : "#64748b",
              }}
            >
              {brand.cagr !== null && brand.cagr >= 0 ? <TrendingUp size={20} /> : brand.cagr !== null ? <TrendingDown size={20} /> : <TrendingUp size={20} />}
            </div>
          </div>
          <div className="stat-value">
            {brand.cagr !== null ? `${brand.cagr.toFixed(1)}%` : 'N/A'}
          </div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            {brand.cagr !== null ? '年复合增长率' : '数据不足'}
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">商品数量</span>
            <div className="stat-icon" style={{ background: "#fce7f3", color: "#ec4899" }}>
              <Package size={20} />
            </div>
          </div>
          <div className="stat-value">{brand.productCount}</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            在售商品总数
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">品牌独立站</span>
            <div className="stat-icon" style={{ background: "#dbeafe", color: "#2563eb" }}>
              <Globe size={20} />
            </div>
          </div>
          <div style={{ marginTop: "12px" }}>
            {brand.website ? (
              <a
                href={brand.website}
                target="_blank"
                rel="noopener noreferrer"
                className="button button-primary"
                style={{ width: "100%" }}
              >
                访问网站
                <ExternalLink size={16} />
              </a>
            ) : (
              <span style={{ color: "var(--text-secondary)" }}>暂无数据</span>
            )}
          </div>
        </div>
      </div>

      {/* Sales Trend Chart */}
      <div className="card" style={{ marginTop: "32px" }}>
        <div className="card-header">
          <h2 className="card-title">品牌销售趋势</h2>
          <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
            品牌销售额历史趋势（{brand.monthlyTrends.length}个月数据）
          </p>
        </div>
        <div className="card-content">
          <ResponsiveContainer width="100%" height={400}>
            <LineChart data={trendData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
              <XAxis dataKey="date" stroke="#64748b" style={{ fontSize: "12px" }} />
              <YAxis
                stroke="#64748b"
                style={{ fontSize: "12px" }}
                tickFormatter={(value) => formatNumber(value)}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: "white",
                  border: "1px solid #e2e8f0",
                  borderRadius: "8px",
                }}
                formatter={(value: number) => formatNumber(value)}
              />
              <Legend />
              <Line
                type="monotone"
                dataKey="revenue"
                stroke="#2563eb"
                strokeWidth={3}
                name="销售额"
                dot={{ fill: "#2563eb", r: 4 }}
                activeDot={{ r: 6 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Social Media Section */}
      <div className="card" style={{ marginTop: "32px" }}>
        <div className="card-header">
          <h2 className="card-title">社交媒体数据</h2>
          <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
            品牌在各社交媒体平台的表现
          </p>
        </div>
        <div className="card-content">
          <div className="stats-grid">
            {brand.youtube?.url && (
              <div className="stat-card">
                <div className="stat-header">
                  <span className="stat-label">YouTube</span>
                  <Youtube size={24} color="#ff0000" />
                </div>
                <div className="stat-value" style={{ fontSize: "24px" }}>
                  {formatVolume(brand.youtube.subscribers || 0)}
                </div>
                <div style={{ fontSize: "14px", color: "var(--text-secondary)", marginBottom: "12px" }}>
                  订阅用户
                </div>
                <a
                  href={brand.youtube.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="button button-secondary"
                  style={{ width: "100%", justifyContent: "center" }}
                >
                  访问频道
                  <ExternalLink size={14} />
                </a>
              </div>
            )}

            {brand.instagram?.url && (
              <div className="stat-card">
                <div className="stat-header">
                  <span className="stat-label">Instagram</span>
                  <Instagram size={24} color="#e4405f" />
                </div>
                <div className="stat-value" style={{ fontSize: "24px" }}>
                  {formatVolume(brand.instagram.followers || 0)}
                </div>
                <div style={{ fontSize: "14px", color: "var(--text-secondary)", marginBottom: "12px" }}>
                  粉丝数量
                </div>
                <a
                  href={brand.instagram.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="button button-secondary"
                  style={{ width: "100%", justifyContent: "center" }}
                >
                  访问主页
                  <ExternalLink size={14} />
                </a>
              </div>
            )}

            {brand.facebook?.url && (
              <div className="stat-card">
                <div className="stat-header">
                  <span className="stat-label">Facebook</span>
                  <Facebook size={24} color="#1877f2" />
                </div>
                <div className="stat-value" style={{ fontSize: "24px" }}>
                  {formatVolume(brand.facebook.followers || 0)}
                </div>
                <div style={{ fontSize: "14px", color: "var(--text-secondary)", marginBottom: "12px" }}>
                  关注数量
                </div>
                <a
                  href={brand.facebook.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="button button-secondary"
                  style={{ width: "100%", justifyContent: "center" }}
                >
                  访问主页
                  <ExternalLink size={14} />
                </a>
              </div>
            )}

            {brand.reddit?.url && (
              <div className="stat-card">
                <div className="stat-header">
                  <span className="stat-label">Reddit</span>
                  <MessageCircle size={24} color="#ff4500" />
                </div>
                <div className="stat-value" style={{ fontSize: "24px" }}>
                  {formatVolume(brand.reddit.posts || 0)}
                </div>
                <div style={{ fontSize: "14px", color: "var(--text-secondary)", marginBottom: "12px" }}>
                  讨论数量
                </div>
                <a
                  href={brand.reddit.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="button button-secondary"
                  style={{ width: "100%", justifyContent: "center" }}
                >
                  访问社区
                  <ExternalLink size={14} />
                </a>
              </div>
            )}
          </div>

          {!brand.youtube?.url &&
            !brand.instagram?.url &&
            !brand.facebook?.url &&
            !brand.reddit?.url && (
              <div style={{ textAlign: "center", padding: "40px", color: "var(--text-secondary)" }}>
                暂无社交媒体数据
              </div>
            )}
        </div>
      </div>

      {/* Products Section */}
      <div className="card" style={{ marginTop: "32px" }}>
        <div className="card-header">
          <h2 className="card-title">品牌商品</h2>
          <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
            该品牌在亚马逊的所有商品列表
          </p>
        </div>
        <div className="card-content">
          {brandProducts.length > 0 ? (
            <div className="product-grid">
              {brandProducts.map((product) => {
                // 获取商品销售趋势数据
                const productSales = market.productSales[product.asin];
                const last6Months = productSales?.monthlySales?.slice(-6) || [];
                const hasSalesData = last6Months.length > 0;
                
                // 计算销售趋势统计
                let salesTrend = 0;
                let maxSales = 0;
                if (last6Months.length >= 2) {
                  const firstMonth = last6Months[0]?.sales || 0;
                  const lastMonth = last6Months[last6Months.length - 1]?.sales || 0;
                  if (firstMonth > 0) {
                    salesTrend = ((lastMonth - firstMonth) / firstMonth) * 100;
                  }
                  maxSales = Math.max(...last6Months.map(m => m.sales || 0));
                }
                
                return (
                  <div key={product.asin} className="product-card">
                    {product.imageUrl ? (
                      <img
                        src={product.imageUrl}
                        alt={product.title}
                        className="product-image"
                        onError={(e) => {
                          (e.target as HTMLImageElement).src = "/placeholder-product.png";
                        }}
                      />
                    ) : (
                      <div
                        className="product-image"
                        style={{
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          background: "var(--background)",
                        }}
                      >
                        <Package size={48} color="var(--text-secondary)" />
                      </div>
                    )}
                    <div className="product-info">
                      <div className="product-title" title={product.title}>
                        {product.title || "No title available"}
                      </div>
                      <div className="product-brand">{product.brand}</div>
                      {product.rating > 0 && (
                        <div className="product-rating">
                          <Star size={14} fill="#f59e0b" />
                          <span style={{ fontWeight: 600 }}>{product.rating.toFixed(1)}</span>
                          <span style={{ color: "var(--text-secondary)" }}>
                            ({product.reviews.toLocaleString()})
                          </span>
                        </div>
                      )}
                      
                      {/* 销售趋势迷你图 */}
                      {hasSalesData && (
                        <div style={{ marginTop: "12px", marginBottom: "8px" }}>
                          <div style={{ 
                            display: "flex", 
                            justifyContent: "space-between", 
                            alignItems: "center",
                            marginBottom: "6px"
                          }}>
                            <span style={{ fontSize: "11px", color: "var(--text-secondary)", fontWeight: 500 }}>
                              销售趋势 (6个月)
                            </span>
                            {salesTrend !== 0 && (
                              <span style={{ 
                                fontSize: "11px", 
                                fontWeight: 600,
                                color: salesTrend > 0 ? "#10b981" : "#ef4444"
                              }}>
                                {salesTrend > 0 ? "↑" : "↓"} {Math.abs(salesTrend).toFixed(0)}%
                              </span>
                            )}
                          </div>
                          <div style={{ 
                            display: "flex", 
                            alignItems: "flex-end", 
                            height: "40px", 
                            gap: "2px",
                            padding: "4px",
                            background: "var(--background)",
                            borderRadius: "4px"
                          }}>
                            {last6Months.map((month, idx) => {
                              const height = maxSales > 0 ? ((month.sales || 0) / maxSales) * 100 : 0;
                              const isLatest = idx === last6Months.length - 1;
                              return (
                                <div
                                  key={idx}
                                  style={{
                                    flex: 1,
                                    height: `${height}%`,
                                    minHeight: height > 0 ? "4px" : "2px",
                                    background: isLatest 
                                      ? salesTrend >= 0 ? "#10b981" : "#ef4444"
                                      : "#2563eb",
                                    borderRadius: "2px 2px 0 0",
                                    opacity: 0.4 + (idx / last6Months.length) * 0.6,
                                    transition: "all 0.2s",
                                  }}
                                  title={`${month.date}: ${formatNumber(month.sales || 0)}`}
                                />
                              );
                            })}
                          </div>
                        </div>
                      )}
                      
                      <div
                        style={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          marginTop: "8px",
                        }}
                      >
                        <div className="product-price">
                          {product.price > 0 ? formatNumber(product.price) : "N/A"}
                        </div>
                        {product.monthlySales > 0 && (
                          <div
                            style={{
                              fontSize: "12px",
                              color: "var(--text-secondary)",
                            }}
                          >
                            ~{product.monthlySales.toLocaleString()}/月
                          </div>
                        )}
                      </div>
                      <a
                        href={`https://www.amazon.com/dp/${product.asin}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="button button-primary"
                        style={{
                          width: "100%",
                          marginTop: "12px",
                          justifyContent: "center",
                          fontSize: "13px",
                          padding: "8px 16px",
                        }}
                        onClick={(e) => e.stopPropagation()}
                      >
                        在亚马逊查看
                        <ExternalLink size={14} />
                      </a>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div style={{ textAlign: "center", padding: "40px", color: "var(--text-secondary)" }}>
              暂无商品数据
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

