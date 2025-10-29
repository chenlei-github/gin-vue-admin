import type { Metadata } from "next";
import "./globals.css";
import { TrendingUp } from "lucide-react";
import Link from "next/link";

export const metadata: Metadata = {
  title: "BrandTrekin Beta - 商业数据分析平台",
  description: "基于线上商业数据分析的品牌卖家数据支持与建议平台",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body>
        <header className="header">
          <div className="container">
            <div className="header-content">
              <Link href="/" className="logo">
                <div className="logo-icon">
                  <TrendingUp size={20} />
                </div>
                <span>BrandTrekin Beta</span>
              </Link>
              <nav>
                <span style={{ fontSize: '14px', color: 'var(--text-secondary)' }}>
                  商业数据分析平台
                </span>
              </nav>
            </div>
          </div>
        </header>
        <main>{children}</main>
      </body>
    </html>
  );
}
