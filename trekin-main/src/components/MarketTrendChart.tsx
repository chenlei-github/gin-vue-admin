"use client";

import { useState, useMemo } from "react";
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
import { MarketDetail, Keyword } from "@/types";
import { X, Plus } from "lucide-react";

interface MarketTrendChartProps {
  market: MarketDetail;
  allMarkets?: { id: string; name: string; metrics: any }[];
}

const CHART_COLORS = [
  "#2563eb", // blue
  "#f59e0b", // orange
  "#10b981", // green
  "#ef4444", // red
  "#8b5cf6", // purple
  "#ec4899", // pink
  "#06b6d4", // cyan
  "#f97316", // deep orange
];

export default function MarketTrendChart({ market, allMarkets = [] }: MarketTrendChartProps) {
  // 默认指标
  const [selectedMetrics, setSelectedMetrics] = useState<string[]>([
    "market-revenue",
    "total-search",
  ]);

  const [showMetricSelector, setShowMetricSelector] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  // 准备可用的指标选项
  const availableMetrics = useMemo(() => {
    const metrics: { id: string; name: string; type: string; axis: string }[] = [];

    // 当前市场销售额
    metrics.push({
      id: "market-revenue",
      name: `${market.name} - 总销售额`,
      type: "revenue",
      axis: "left",
    });

    // 总搜索量
    metrics.push({
      id: "total-search",
      name: "总搜索量",
      type: "search",
      axis: "right",
    });

    // Google搜索量
    metrics.push({
      id: "google-search",
      name: "Google搜索量",
      type: "search",
      axis: "right",
    });

    // Amazon搜索量
    metrics.push({
      id: "amazon-search",
      name: "Amazon搜索量",
      type: "search",
      axis: "right",
    });

    // 各个关键词 - 分别处理Google和Amazon，显示所有关键词
    const googleKeywords = market.keywords.filter((k) => k.source === "google");
    const amazonKeywords = market.keywords.filter((k) => k.source === "amazon");
    
    googleKeywords.forEach((keyword, idx) => {
      const displayName = keyword.keyword.length > 40 
        ? keyword.keyword.substring(0, 40) + "..." 
        : keyword.keyword;
      metrics.push({
        id: `keyword-google-${idx}`,
        name: `${displayName} (Google)`,
        type: "search",
        axis: "right",
      });
    });
    
    amazonKeywords.forEach((keyword, idx) => {
      // 清理Amazon关键词名称，去掉"History-"前缀和"-US"后缀
      let cleanName = keyword.keyword
        .replace(/^History-/, '')
        .replace(/-US$/, '');
      
      const displayName = cleanName.length > 40 
        ? cleanName.substring(0, 40) + "..." 
        : cleanName;
      metrics.push({
        id: `keyword-amazon-${idx}`,
        name: `${displayName} (Amazon)`,
        type: "search",
        axis: "right",
      });
    });

    // 所有品牌 - 按销售额排序
    market.brands.forEach((brand) => {
      metrics.push({
        id: `brand-${brand.brand}`,
        name: `品牌: ${brand.brand}`,
        type: "revenue",
        axis: "left",
      });
    });

    // 其他市场 - 用于对比
    allMarkets.forEach((otherMarket) => {
      if (otherMarket.id !== market.id) {
        metrics.push({
          id: `compare-market-${otherMarket.id}`,
          name: `对比: ${otherMarket.name}`,
          type: "revenue",
          axis: "left",
        });
      }
    });

    return metrics;
  }, [market, allMarkets]);

  // 准备图表数据
  const chartData = useMemo(() => {
    const data: any[] = [];

    market.metrics.monthlyTrends.forEach((trend) => {
      const dataPoint: any = { date: trend.date };

      selectedMetrics.forEach((metricId) => {
        if (metricId === "market-revenue") {
          dataPoint[metricId] = trend.revenue || 0;
        } else if (metricId === "total-search") {
          const keywordTrend = market.metrics.keywordMonthlyTrends.find(
            (k) => k.date === trend.date
          );
          dataPoint[metricId] = keywordTrend?.volume || 0;
        } else if (metricId === "google-search") {
          // 计算Google搜索量总和
          let googleVolume = 0;
          market.keywords
            .filter((k) => k.source === "google")
            .forEach((keyword) => {
              const monthData = keyword.monthlyVolume.find((m) => m.date === trend.date);
              googleVolume += monthData?.volume || 0;
            });
          dataPoint[metricId] = googleVolume;
        } else if (metricId === "amazon-search") {
          // 计算Amazon搜索量总和
          let amazonVolume = 0;
          market.keywords
            .filter((k) => k.source === "amazon")
            .forEach((keyword) => {
              const monthData = keyword.monthlyVolume.find((m) => m.date === trend.date);
              amazonVolume += monthData?.volume || 0;
            });
          dataPoint[metricId] = amazonVolume;
        } else if (metricId.startsWith("keyword-google-")) {
          // Google关键词
          const idx = parseInt(metricId.replace("keyword-google-", ""));
          const googleKeywords = market.keywords.filter((k) => k.source === "google");
          const keyword = googleKeywords[idx];
          if (keyword) {
            const monthData = keyword.monthlyVolume.find((m) => m.date === trend.date);
            dataPoint[metricId] = monthData?.volume || 0;
          }
        } else if (metricId.startsWith("keyword-amazon-")) {
          // Amazon关键词
          const idx = parseInt(metricId.replace("keyword-amazon-", ""));
          const amazonKeywords = market.keywords.filter((k) => k.source === "amazon");
          const keyword = amazonKeywords[idx];
          if (keyword) {
            const monthData = keyword.monthlyVolume.find((m) => m.date === trend.date);
            dataPoint[metricId] = monthData?.volume || 0;
          }
        } else if (metricId.startsWith("brand-")) {
          // 品牌销售额
          const brandName = metricId.replace("brand-", "");
          const brand = market.brands.find((b) => b.brand === brandName);
          if (brand) {
            const monthData = brand.monthlyTrends.find((m) => m.date === trend.date);
            dataPoint[metricId] = monthData?.revenue || 0;
          }
        } else if (metricId.startsWith("compare-market-")) {
          // 其他市场对比 - 使用摘要数据中的趋势
          const otherMarketId = metricId.replace("compare-market-", "");
          const otherMarket = allMarkets.find((m) => m.id === otherMarketId);
          if (otherMarket && otherMarket.metrics?.monthlyTrends) {
            const monthData = otherMarket.metrics.monthlyTrends.find((m: any) => m.date === trend.date);
            dataPoint[metricId] = monthData?.revenue || 0;
          }
        }
      });

      data.push(dataPoint);
    });

    return data;
  }, [market, selectedMetrics, allMarkets]);

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

  const toggleMetric = (metricId: string) => {
    if (selectedMetrics.includes(metricId)) {
      setSelectedMetrics(selectedMetrics.filter((id) => id !== metricId));
    } else {
      setSelectedMetrics([...selectedMetrics, metricId]);
    }
  };

  const selectedMetricDetails = selectedMetrics
    .map((id) => availableMetrics.find((m) => m.id === id))
    .filter(Boolean);

  // 过滤指标
  const filteredMetrics = useMemo(() => {
    if (!searchQuery) return availableMetrics;
    const query = searchQuery.toLowerCase();
    return availableMetrics.filter((m) =>
      m.name.toLowerCase().includes(query)
    );
  }, [availableMetrics, searchQuery]);

  // 按类型分组
  const groupedMetrics = useMemo(() => {
    const groups: { [key: string]: typeof availableMetrics } = {
      market: [],
      brand: [],
      searchTotal: [],
      keyword: [],
      compare: [],
    };

    filteredMetrics.forEach((metric) => {
      if (metric.id === "market-revenue") {
        groups.market.push(metric);
      } else if (metric.id.startsWith("brand-")) {
        groups.brand.push(metric);
      } else if (
        metric.id === "total-search" ||
        metric.id === "google-search" ||
        metric.id === "amazon-search"
      ) {
        groups.searchTotal.push(metric);
      } else if (metric.id.startsWith("keyword-")) {
        groups.keyword.push(metric);
      } else if (metric.id.startsWith("compare-")) {
        groups.compare.push(metric);
      }
    });

    return groups;
  }, [filteredMetrics]);

  return (
    <div className="card">
      <div className="card-header">
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "start" }}>
          <div>
            <h2 className="card-title">市场趋势分析</h2>
            <p style={{ fontSize: "14px", color: "var(--text-secondary)", marginTop: "4px" }}>
              自定义指标进行交叉分析（{market.metrics.monthlyTrends.length}个月数据）
            </p>
          </div>
          <button
            onClick={() => setShowMetricSelector(!showMetricSelector)}
            className="button button-primary"
            style={{ fontSize: "13px", padding: "8px 16px" }}
          >
            <Plus size={16} />
            添加指标
          </button>
        </div>

        {/* 已选择的指标 */}
        <div style={{ marginTop: "16px", display: "flex", flexWrap: "wrap", gap: "8px" }}>
          {selectedMetricDetails.map((metric, idx) => (
            <div
              key={metric!.id}
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: "6px",
                padding: "6px 12px",
                background: "var(--background)",
                border: "1px solid var(--border)",
                borderRadius: "8px",
                fontSize: "13px",
              }}
            >
              <div
                style={{
                  width: "12px",
                  height: "12px",
                  borderRadius: "50%",
                  background: CHART_COLORS[idx % CHART_COLORS.length],
                }}
              />
              <span>{metric!.name}</span>
              {selectedMetrics.length > 1 && (
                <button
                  onClick={() => toggleMetric(metric!.id)}
                  style={{
                    background: "none",
                    border: "none",
                    cursor: "pointer",
                    padding: "2px",
                    display: "flex",
                    color: "var(--text-secondary)",
                  }}
                >
                  <X size={14} />
                </button>
              )}
            </div>
          ))}
        </div>

        {/* 指标选择器 */}
        {showMetricSelector && (
          <div
            style={{
              marginTop: "16px",
              padding: "16px",
              background: "var(--background)",
              border: "1px solid var(--border)",
              borderRadius: "8px",
              maxHeight: "500px",
              overflowY: "auto",
            }}
          >
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "12px" }}>
              <h3 style={{ fontSize: "14px", fontWeight: 600 }}>选择指标</h3>
              <button
                onClick={() => {
                  setShowMetricSelector(false);
                  setSearchQuery("");
                }}
                className="button button-secondary"
                style={{ fontSize: "12px", padding: "4px 12px" }}
              >
                关闭
              </button>
            </div>

            {/* 搜索框 */}
            <div style={{ marginBottom: "16px" }}>
              <input
                type="text"
                placeholder="搜索指标（品牌名称、关键词等）..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{
                  width: "100%",
                  padding: "8px 12px",
                  border: "1px solid var(--border)",
                  borderRadius: "6px",
                  fontSize: "13px",
                }}
              />
            </div>

            {/* 市场整体 */}
            {groupedMetrics.market.length > 0 && (
              <div style={{ marginBottom: "16px" }}>
                <h4
                  style={{
                    fontSize: "13px",
                    fontWeight: 600,
                    color: "var(--text-secondary)",
                    marginBottom: "8px",
                  }}
                >
                  市场整体
                </h4>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {groupedMetrics.market.map((metric) => (
                    <label
                      key={metric.id}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        padding: "8px",
                        borderRadius: "6px",
                        cursor: "pointer",
                        fontSize: "13px",
                      }}
                      className="metric-option"
                    >
                      <input
                        type="checkbox"
                        checked={selectedMetrics.includes(metric.id)}
                        onChange={() => toggleMetric(metric.id)}
                        style={{ cursor: "pointer" }}
                      />
                      <span>{metric.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {/* 品牌销售 */}
            {groupedMetrics.brand.length > 0 && (
              <div style={{ marginBottom: "16px" }}>
                <h4
                  style={{
                    fontSize: "13px",
                    fontWeight: 600,
                    color: "var(--text-secondary)",
                    marginBottom: "8px",
                  }}
                >
                  品牌销售额 ({groupedMetrics.brand.length} 个品牌)
                </h4>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {groupedMetrics.brand.map((metric) => (
                    <label
                      key={metric.id}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        padding: "8px",
                        borderRadius: "6px",
                        cursor: "pointer",
                        fontSize: "13px",
                      }}
                      className="metric-option"
                    >
                      <input
                        type="checkbox"
                        checked={selectedMetrics.includes(metric.id)}
                        onChange={() => toggleMetric(metric.id)}
                        style={{ cursor: "pointer" }}
                      />
                      <span>{metric.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {/* 搜索量汇总 */}
            {groupedMetrics.searchTotal.length > 0 && (
              <div style={{ marginBottom: "16px" }}>
                <h4
                  style={{
                    fontSize: "13px",
                    fontWeight: 600,
                    color: "var(--text-secondary)",
                    marginBottom: "8px",
                  }}
                >
                  搜索量汇总
                </h4>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {groupedMetrics.searchTotal.map((metric) => (
                    <label
                      key={metric.id}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        padding: "8px",
                        borderRadius: "6px",
                        cursor: "pointer",
                        fontSize: "13px",
                      }}
                      className="metric-option"
                    >
                      <input
                        type="checkbox"
                        checked={selectedMetrics.includes(metric.id)}
                        onChange={() => toggleMetric(metric.id)}
                        style={{ cursor: "pointer" }}
                      />
                      <span>{metric.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {/* 关键词搜索量 */}
            {groupedMetrics.keyword.length > 0 && (
              <div style={{ marginBottom: "16px" }}>
                <h4
                  style={{
                    fontSize: "13px",
                    fontWeight: 600,
                    color: "var(--text-secondary)",
                    marginBottom: "8px",
                  }}
                >
                  关键词搜索量 ({groupedMetrics.keyword.length} 个关键词)
                </h4>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {groupedMetrics.keyword.map((metric) => (
                    <label
                      key={metric.id}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        padding: "8px",
                        borderRadius: "6px",
                        cursor: "pointer",
                        fontSize: "13px",
                      }}
                      className="metric-option"
                    >
                      <input
                        type="checkbox"
                        checked={selectedMetrics.includes(metric.id)}
                        onChange={() => toggleMetric(metric.id)}
                        style={{ cursor: "pointer" }}
                      />
                      <span>{metric.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {/* 市场对比 */}
            {groupedMetrics.compare.length > 0 && (
              <div style={{ marginBottom: "16px" }}>
                <h4
                  style={{
                    fontSize: "13px",
                    fontWeight: 600,
                    color: "var(--text-secondary)",
                    marginBottom: "8px",
                  }}
                >
                  市场对比
                </h4>
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {groupedMetrics.compare.map((metric) => (
                    <label
                      key={metric.id}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        padding: "8px",
                        borderRadius: "6px",
                        cursor: "pointer",
                        fontSize: "13px",
                      }}
                      className="metric-option"
                    >
                      <input
                        type="checkbox"
                        checked={selectedMetrics.includes(metric.id)}
                        onChange={() => toggleMetric(metric.id)}
                        style={{ cursor: "pointer" }}
                      />
                      <span>{metric.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            )}

            {filteredMetrics.length === 0 && (
              <div
                style={{
                  textAlign: "center",
                  padding: "24px",
                  color: "var(--text-secondary)",
                }}
              >
                未找到匹配的指标
              </div>
            )}
          </div>
        )}
      </div>

      <div className="card-content">
        <ResponsiveContainer width="100%" height={450}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
            <XAxis dataKey="date" stroke="#64748b" style={{ fontSize: "12px" }} />
            
            {/* 左轴：销售额 */}
            {selectedMetricDetails.some((m) => m?.type === "revenue") && (
              <YAxis
                yAxisId="left"
                stroke="#2563eb"
                style={{ fontSize: "12px" }}
                tickFormatter={(value) => formatNumber(value)}
              />
            )}
            
            {/* 右轴：搜索量 */}
            {selectedMetricDetails.some((m) => m?.type === "search") && (
              <YAxis
                yAxisId="right"
                orientation="right"
                stroke="#f59e0b"
                style={{ fontSize: "12px" }}
                tickFormatter={(value) => formatVolume(value)}
              />
            )}
            
            <Tooltip
              contentStyle={{
                backgroundColor: "white",
                border: "1px solid #e2e8f0",
                borderRadius: "8px",
              }}
              formatter={(value: number, name: string) => {
                const metric = availableMetrics.find((m) => m.id === name);
                if (metric?.type === "revenue") return formatNumber(value);
                if (metric?.type === "search") return formatVolume(value);
                return value;
              }}
              labelFormatter={(label) => `时间: ${label}`}
            />
            
            <Legend
              formatter={(value) => {
                const metric = availableMetrics.find((m) => m.id === value);
                return metric?.name || value;
              }}
            />

            {selectedMetricDetails.map((metric, idx) => {
              if (!metric) return null;
              return (
                <Line
                  key={metric.id}
                  yAxisId={metric.type === "revenue" ? "left" : "right"}
                  type="monotone"
                  dataKey={metric.id}
                  stroke={CHART_COLORS[idx % CHART_COLORS.length]}
                  strokeWidth={2}
                  name={metric.id}
                  dot={false}
                  connectNulls
                />
              );
            })}
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

