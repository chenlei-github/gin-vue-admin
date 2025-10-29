"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import dynamic from "next/dynamic";
import {
  TrendingUp,
  TrendingDown,
  Building2,
  Search,
  ArrowLeft,
  ChevronRight,
  Youtube,
  Instagram,
  Facebook,
  Globe,
  MessageCircle,
} from "lucide-react";
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { MarketDetail, BrandMetrics, MarketSummary } from "@/types";

// 动态导入趋势图组件
const MarketTrendChart = dynamic(() => import("@/components/MarketTrendChart"), {
  ssr: false,
});

const COLORS = [
  "#2563eb",
  "#7c3aed",
  "#ec4899",
  "#f59e0b",
  "#10b981",
  "#06b6d4",
  "#8b5cf6",
  "#f97316",
];

export default function MarketDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [market, setMarket] = useState<MarketDetail | null>(null);
  const [allMarkets, setAllMarkets] = useState<MarketSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 加载当前市场和所有市场数据
    Promise.all([
      params.id ? fetch(`/data/${params.id}.json`).then((res) => res.json()) : null,
      fetch(`/data/markets.json`).then((res) => res.json()),
    ])
      .then(([marketData, marketsData]) => {
        if (marketData) setMarket(marketData);
        setAllMarkets(marketsData || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error loading data:", err);
        setLoading(false);
      });
  }, [params.id]);

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

  if (!market) {
    return (
      <div className="container" style={{ padding: "48px 24px", textAlign: "center" }}>
        <h1>市场未找到</h1>
        <Link href="/">
          <button className="button button-primary" style={{ marginTop: "24px" }}>
            返回首页
          </button>
        </Link>
      </div>
    );
  }

  // Top brands for pie chart
  const topBrands = market.brands.slice(0, 8);
  const pieData = topBrands.map((brand) => ({
    name: brand.brand,
    value: brand.totalRevenue,
  }));

  // Brand bar chart data
  const barData = topBrands.map((brand) => ({
    name: brand.brand.length > 15 ? brand.brand.substring(0, 15) + "..." : brand.brand,
    revenue: brand.totalRevenue,
  }));

  return (
    <div className="container" style={{ padding: "48px 24px" }}>
      {/* Breadcrumb */}
      <div className="breadcrumb">
        <Link href="/">市场分析</Link>
        <ChevronRight size={16} />
        <span>{market.name}</span>
      </div>

      {/* Header */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "start", marginBottom: "32px" }}>
        <div>
          <h1 className="page-title">{market.name}</h1>
          <p className="page-subtitle">
            分析该细分市场的竞争情况、品牌表现和市场趋势
          </p>
        </div>
        <Link href="/">
          <button className="button button-secondary">
            <ArrowLeft size={16} />
            返回
          </button>
        </Link>
      </div>

      {/* Key Metrics */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">市场规模</span>
            <div className="stat-icon" style={{ background: "#dcfce7", color: "#10b981" }}>
              <TrendingUp size={20} />
            </div>
          </div>
          <div className="stat-value">{formatNumber(market.metrics.totalRevenue)}</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            年度市场规模预估
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">市场增速</span>
            <div
              className="stat-icon"
              style={{
                background: market.metrics.cagr !== null && market.metrics.cagr >= 0 ? "#dcfce7" : market.metrics.cagr !== null ? "#fee2e2" : "#f1f5f9",
                color: market.metrics.cagr !== null && market.metrics.cagr >= 0 ? "#10b981" : market.metrics.cagr !== null ? "#ef4444" : "#64748b",
              }}
            >
              {market.metrics.cagr !== null && market.metrics.cagr >= 0 ? (
                <TrendingUp size={20} />
              ) : market.metrics.cagr !== null ? (
                <TrendingDown size={20} />
              ) : (
                <TrendingUp size={20} />
              )}
            </div>
          </div>
          <div className="stat-value">
            {market.metrics.cagr !== null ? `${market.metrics.cagr.toFixed(1)}%` : 'N/A'}
          </div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            {market.metrics.cagr !== null ? '年复合增长率' : '数据不足'}
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">市场声量</span>
            <div className="stat-icon" style={{ background: "#fef3c7", color: "#f59e0b" }}>
              <Search size={20} />
            </div>
          </div>
          <div className="stat-value">{formatVolume(market.metrics.searchVolume)}</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            月度搜索量总和
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">品牌数量</span>
            <div className="stat-icon" style={{ background: "#dbeafe", color: "#2563eb" }}>
              <Building2 size={20} />
            </div>
          </div>
          <div className="stat-value">{market.metrics.brandCount}</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
            活跃品牌总数
          </div>
        </div>
      </div>

      {/* Market Trends Chart - Interactive */}
      <div style={{ marginTop: "32px" }}>
        <MarketTrendChart market={market} allMarkets={allMarkets} />
      </div>

      {/* Brand Distribution */}
      <div className="grid grid-cols-2" style={{ marginTop: "32px" }}>
        <div className="card">
          <div className="card-header">
            <h2 className="card-title">品牌销售额占比</h2>
            <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
              Top 8 品牌市场份额分布
            </p>
          </div>
          <div className="card-content">
            <ResponsiveContainer width="100%" height={350}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) =>
                    `${name.length > 12 ? name.substring(0, 12) + "..." : name} ${(percent * 100).toFixed(0)}%`
                  }
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {pieData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value: number) => formatNumber(value)}
                  contentStyle={{
                    backgroundColor: "white",
                    border: "1px solid #e2e8f0",
                    borderRadius: "8px",
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="card">
          <div className="card-header">
            <h2 className="card-title">品牌销售额排名</h2>
            <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
              Top 8 品牌年度销售额对比
            </p>
          </div>
          <div className="card-content">
            <ResponsiveContainer width="100%" height={350}>
              <BarChart data={barData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                <XAxis
                  dataKey="name"
                  stroke="#64748b"
                  style={{ fontSize: "12px" }}
                  angle={-45}
                  textAnchor="end"
                  height={100}
                />
                <YAxis
                  stroke="#64748b"
                  style={{ fontSize: "12px" }}
                  tickFormatter={(value) => formatNumber(value)}
                />
                <Tooltip
                  formatter={(value: number) => formatNumber(value)}
                  contentStyle={{
                    backgroundColor: "white",
                    border: "1px solid #e2e8f0",
                    borderRadius: "8px",
                  }}
                />
                <Bar dataKey="revenue" fill="#2563eb" radius={[8, 8, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Brand List */}
      <div className="card" style={{ marginTop: "32px" }}>
        <div className="card-header">
          <h2 className="card-title">品牌列表</h2>
          <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
            点击品牌查看详细信息和产品列表
          </p>
        </div>
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>品牌名称</th>
                <th>品牌独立站</th>
                <th>品牌规模</th>
                <th>品牌增速</th>
                <th>社交媒体热度</th>
                <th style={{ textAlign: "right" }}>操作</th>
              </tr>
            </thead>
            <tbody>
              {market.brands.map((brand) => (
                <tr key={brand.brand}>
                  <td>
                    <div style={{ fontWeight: 600 }}>{brand.brand}</div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      {brand.productCount} 个商品
                    </div>
                  </td>
                  <td>
                    {brand.website ? (
                      <a
                        href={brand.website}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="social-link"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <Globe size={14} />
                        访问网站
                      </a>
                    ) : (
                      <span style={{ color: "var(--text-secondary)" }}>-</span>
                    )}
                  </td>
                  <td>
                    <div style={{ fontWeight: 600 }}>
                      {formatNumber(brand.totalRevenue)}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      年度预估
                    </div>
                  </td>
                  <td>
                    {brand.cagr !== null ? (
                      <div
                        className={`stat-change ${
                          brand.cagr >= 0 ? "positive" : "negative"
                        }`}
                      >
                        {brand.cagr >= 0 ? (
                          <TrendingUp size={16} />
                        ) : (
                          <TrendingDown size={16} />
                        )}
                        <span style={{ fontWeight: 600 }}>{brand.cagr.toFixed(1)}%</span>
                      </div>
                    ) : (
                      <div style={{ color: "var(--text-secondary)", fontSize: "14px" }}>
                        N/A
                      </div>
                    )}
                  </td>
                  <td>
                    <div className="social-links" onClick={(e) => e.stopPropagation()}>
                      {brand.youtube?.url && (
                        <a
                          href={brand.youtube.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="social-link"
                          title={`${brand.youtube.subscribers?.toLocaleString()} subscribers`}
                        >
                          <Youtube size={14} color="#ff0000" />
                          {formatVolume(brand.youtube.subscribers || 0)}
                        </a>
                      )}
                      {brand.instagram?.url && (
                        <a
                          href={brand.instagram.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="social-link"
                          title={`${brand.instagram.followers?.toLocaleString()} followers`}
                        >
                          <Instagram size={14} color="#e4405f" />
                          {formatVolume(brand.instagram.followers || 0)}
                        </a>
                      )}
                      {brand.facebook?.url && (
                        <a
                          href={brand.facebook.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="social-link"
                          title={`${brand.facebook.followers?.toLocaleString()} followers`}
                        >
                          <Facebook size={14} color="#1877f2" />
                          {formatVolume(brand.facebook.followers || 0)}
                        </a>
                      )}
                      {brand.reddit?.url && (
                        <a
                          href={brand.reddit.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="social-link"
                          title={`${brand.reddit.posts?.toLocaleString()} posts`}
                        >
                          <MessageCircle size={14} color="#ff4500" />
                          {formatVolume(brand.reddit.posts || 0)}
                        </a>
                      )}
                    </div>
                  </td>
                  <td style={{ textAlign: "right" }}>
                    <Link href={`/market/${params.id}/brand/${encodeURIComponent(brand.brand)}`}>
                      <button className="button button-primary">查看详情</button>
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

