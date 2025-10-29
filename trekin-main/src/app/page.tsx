"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { 
  TrendingUp, 
  TrendingDown, 
  Package, 
  Search,
  Building2,
  ArrowRight
} from "lucide-react";
import { MarketSummary } from "@/types";

export default function Home() {
  const [markets, setMarkets] = useState<MarketSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("/data/markets.json")
      .then((res) => res.json())
      .then((data) => {
        setMarkets(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error loading markets:", err);
        setLoading(false);
      });
  }, []);

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

  return (
    <div className="container" style={{ padding: "48px 24px" }}>
      <div style={{ marginBottom: "48px" }}>
        <h1 className="page-title">市场分析</h1>
        <p className="page-subtitle">
          探索和分析不同市场的竞争情况，了解市场趋势和品牌表现
        </p>
      </div>

      {/* Overview Stats */}
      <div className="stats-grid" style={{ marginBottom: "48px" }}>
        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">追踪市场数</span>
            <div className="stat-icon" style={{ background: "#dbeafe", color: "#2563eb" }}>
              <Building2 size={20} />
            </div>
          </div>
          <div className="stat-value">{markets.length}</div>
          <div className="badge badge-success">已激活市场</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">总市场规模</span>
            <div className="stat-icon" style={{ background: "#dcfce7", color: "#10b981" }}>
              <TrendingUp size={20} />
            </div>
          </div>
          <div className="stat-value">
            {formatNumber(markets.reduce((sum, m) => sum + m.metrics.totalRevenue, 0))}
          </div>
          <div className="stat-change positive">
            <TrendingUp size={16} />
            <span>平均增长</span>
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">商品总数</span>
            <div className="stat-icon" style={{ background: "#fce7f3", color: "#ec4899" }}>
              <Package size={20} />
            </div>
          </div>
          <div className="stat-value">
            {markets.reduce((sum, m) => sum + m.metrics.totalProducts, 0).toLocaleString()}
          </div>
          <div className="badge badge-info">跨所有市场</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <span className="stat-label">搜索量</span>
            <div className="stat-icon" style={{ background: "#fef3c7", color: "#f59e0b" }}>
              <Search size={20} />
            </div>
          </div>
          <div className="stat-value">
            {formatVolume(markets.reduce((sum, m) => sum + m.metrics.searchVolume, 0))}
          </div>
          <div className="badge badge-warning">月度</div>
        </div>
      </div>

      {/* Markets Table */}
      <div className="card">
        <div className="card-header">
          <h2 className="card-title">市场列表</h2>
        </div>
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>市场名称</th>
                <th>市场规模</th>
                <th>市场增速</th>
                <th>市场声量</th>
                <th>品牌数量</th>
                <th>商品数量</th>
                <th style={{ textAlign: "right" }}>操作</th>
              </tr>
            </thead>
            <tbody>
              {markets.map((market) => (
                <tr key={market.id}>
                  <td>
                    <div style={{ fontWeight: 600, marginBottom: "4px" }}>
                      {market.name}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      亚马逊美国市场
                    </div>
                  </td>
                  <td>
                    <div style={{ fontWeight: 600 }}>
                      {formatNumber(market.metrics.totalRevenue)}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      年度预估
                    </div>
                  </td>
                  <td>
                    {market.metrics.cagr !== null ? (
                      <>
                        <div
                          className={`stat-change ${
                            market.metrics.cagr >= 0 ? "positive" : "negative"
                          }`}
                        >
                          {market.metrics.cagr >= 0 ? (
                            <TrendingUp size={16} />
                          ) : (
                            <TrendingDown size={16} />
                          )}
                          <span style={{ fontWeight: 600 }}>
                            {market.metrics.cagr.toFixed(1)}%
                          </span>
                        </div>
                        <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                          年复合增长率
                        </div>
                      </>
                    ) : (
                      <div style={{ color: "var(--text-secondary)" }}>
                        <div style={{ fontWeight: 600 }}>N/A</div>
                        <div style={{ fontSize: "12px" }}>数据不足</div>
                      </div>
                    )}
                  </td>
                  <td>
                    <div style={{ fontWeight: 600 }}>
                      {formatVolume(market.metrics.searchVolume)}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      月搜索量
                    </div>
                  </td>
                  <td>
                    <div style={{ fontWeight: 600 }}>
                      {market.metrics.brandCount}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      活跃品牌
                    </div>
                  </td>
                  <td>
                    <div style={{ fontWeight: 600 }}>
                      {market.metrics.totalProducts.toLocaleString()}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                      在售商品
                    </div>
                  </td>
                  <td style={{ textAlign: "right" }}>
                    <Link href={`/market/${market.id}`}>
                      <button className="button button-primary">
                        查看详情
                        <ArrowRight size={16} />
                      </button>
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Market Trends Overview */}
      <div style={{ marginTop: "32px" }}>
        <div className="grid grid-cols-2">
          {markets.map((market) => {
            const recentTrends = market.metrics.monthlyTrends.slice(-12);
            const firstMonth = recentTrends[0]?.revenue || 0;
            const lastMonth = recentTrends[recentTrends.length - 1]?.revenue || 0;
            const trendPercentage = firstMonth > 0 
              ? ((lastMonth - firstMonth) / firstMonth * 100).toFixed(1)
              : "0.0";

            return (
              <Link key={market.id} href={`/market/${market.id}`}>
                <div className="card" style={{ cursor: "pointer" }}>
                  <div className="card-content">
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "start", marginBottom: "16px" }}>
                      <div>
                        <h3 style={{ fontSize: "18px", fontWeight: 600, marginBottom: "4px" }}>
                          {market.name}
                        </h3>
                        <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>
                          过去12个月趋势
                        </p>
                      </div>
                      <div className={`stat-change ${parseFloat(trendPercentage) >= 0 ? "positive" : "negative"}`}>
                        {parseFloat(trendPercentage) >= 0 ? <TrendingUp size={20} /> : <TrendingDown size={20} />}
                        <span style={{ fontWeight: 600 }}>{trendPercentage}%</span>
                      </div>
                    </div>
                    
                    <div style={{ display: "flex", alignItems: "end", height: "100px", gap: "4px" }}>
                      {recentTrends.map((trend, idx) => {
                        const maxRevenue = Math.max(...recentTrends.map(t => t.revenue || 0));
                        const height = ((trend.revenue || 0) / maxRevenue) * 100;
                        
                        return (
                          <div
                            key={idx}
                            style={{
                              flex: 1,
                              height: `${height}%`,
                              background: "var(--primary-color)",
                              borderRadius: "4px 4px 0 0",
                              opacity: 0.3 + (idx / recentTrends.length) * 0.7,
                            }}
                          />
                        );
                      })}
                    </div>
                  </div>
                </div>
              </Link>
            );
          })}
        </div>
      </div>
    </div>
  );
}
