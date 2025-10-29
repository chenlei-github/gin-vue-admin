const XLSX = require('xlsx');
const fs = require('fs');
const path = require('path');

// Market configurations
const markets = [
  { id: 'cnc-router', name: 'CNC Router Machine', folder: 'CNCRouter' },
  { id: 'laser-engraver', name: 'Laser Engraver', folder: 'LaserEngraver' },
  { id: 'thermal-camera', name: 'Thermal Camera', folder: 'ThermalCamera' }
];

// Parse Product-US.xlsx
function parseProducts(marketFolder) {
  const filePath = path.join(__dirname, '../data', marketFolder, 'Product-US.xlsx');
  const workbook = XLSX.readFile(filePath);
  const sheetName = workbook.SheetNames[0];
  const worksheet = workbook.Sheets[sheetName];
  
  // 使用header: 1来获取原始数组，因为第一行是标题，第二行才是列名
  const rawData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
  
  if (rawData.length < 3) return [];
  
  // 第二行是列标题
  const headers = rawData[1];
  
  // 找到各列的索引
  const asinIndex = headers.findIndex(h => h === 'ASIN' || h === 'asin');
  const titleIndex = headers.findIndex(h => h === '商品标题' || h === 'Product Title' || h === 'Title');
  const brandIndex = headers.findIndex(h => h === '品牌' || h === 'Brand');
  const priceIndex = headers.findIndex(h => h === '价格($)' || h === 'Price' || h === '价格');
  const ratingIndex = headers.findIndex(h => h === '评分' || h === 'Rating');
  const reviewsIndex = headers.findIndex(h => h === '评分数' || h === 'Review Count' || h === 'Reviews');
  const salesIndex = headers.findIndex(h => h === '月销量' || h === 'Monthly Sales');
  const imageIndex = headers.findIndex(h => h === '商品主图' || h === 'Image' || h === '图片');
  
  const products = [];
  
  // 从第三行开始处理数据（跳过标题行和列标题行）
  for (let i = 2; i < rawData.length; i++) {
    const row = rawData[i];
    
    const asin = row[asinIndex];
    if (!asin) continue;
    
    products.push({
      asin: asin,
      title: row[titleIndex] || '',
      brand: row[brandIndex] || 'Unknown',
      price: parseFloat(row[priceIndex]) || 0,
      rating: parseFloat(row[ratingIndex]) || 0,
      reviews: parseInt(row[reviewsIndex]) || 0,
      monthlySales: parseInt(row[salesIndex]) || 0,
      imageUrl: row[imageIndex] || ''
    });
  }
  
  return products;
}

// Parse product-US-sales.xlsx (横向数据结构)
function parseProductSales(marketFolder) {
  const filePath = path.join(__dirname, '../data', marketFolder, 'product-US-sales.xlsx');
  const workbook = XLSX.readFile(filePath);
  
  const salesData = {};
  
  // 处理销量和销售额两个sheet
  workbook.SheetNames.forEach(sheetName => {
    if (sheetName === 'Notes') return;
    
    const worksheet = workbook.Sheets[sheetName];
    const rawData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
    
    if (rawData.length < 2) return;
    
    const headers = rawData[0];
    const asinIndex = headers.findIndex(h => h && (h === 'ASIN' || h === 'asin'));
    
    if (asinIndex === -1) return;
    
    // 找到所有日期列（格式：2025-10 或 2025-10($)）
    const dateColumns = [];
    headers.forEach((header, index) => {
      if (header && typeof header === 'string') {
        // 匹配 YYYY-MM 或 YYYY-MM($) 格式
        const match = header.match(/^(\d{4}-\d{2})/);
        if (match) {
          dateColumns.push({ index, date: match[1] });
        }
      }
    });
    
    // 处理每一行数据
    for (let i = 1; i < rawData.length; i++) {
      const row = rawData[i];
      const asin = row[asinIndex];
      
      if (!asin) continue;
      
      if (!salesData[asin]) {
        salesData[asin] = { asin, monthlySales: [] };
      }
      
      // 提取每个月的数据
      dateColumns.forEach(({ index, date }) => {
        const value = parseFloat(row[index]) || 0;
        
        // 查找是否已有该月的数据
        let monthData = salesData[asin].monthlySales.find(m => m.date === date);
        
        if (!monthData) {
          monthData = { date, sales: 0, units: 0 };
          salesData[asin].monthlySales.push(monthData);
        }
        
        // 根据sheet名称决定是销售额还是销量
        if (sheetName.includes('销售额') || sheetName.toLowerCase().includes('revenue') || sheetName.toLowerCase().includes('sales')) {
          monthData.sales = value;
        } else if (sheetName.includes('销量') || sheetName.toLowerCase().includes('unit') || sheetName.includes('月销量')) {
          monthData.units = value;
        } else {
          // 如果无法确定，默认作为销售额
          monthData.sales = value;
        }
      });
    }
  });
  
  // 对每个产品的月度数据按日期排序
  Object.values(salesData).forEach(product => {
    product.monthlySales.sort((a, b) => a.date.localeCompare(b.date));
  });
  
  return salesData;
}

// Parse GKW.csv (Google Keywords) - 横向数据结构
function parseGoogleKeywords(marketFolder) {
  const filePath = path.join(__dirname, '../data', marketFolder, 'GKW.csv');
  const workbook = XLSX.readFile(filePath);
  const sheetName = workbook.SheetNames[0];
  const worksheet = workbook.Sheets[sheetName];
  const rawData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
  
  if (rawData.length < 2) return [];
  
  const headers = rawData[0];
  const keywords = [];
  
  // 找到所有 "Searches: Month Year" 格式的列
  const searchColumns = [];
  headers.forEach((header, index) => {
    if (header && typeof header === 'string' && header.startsWith('Searches: ')) {
      const dateStr = header.replace('Searches: ', '');
      searchColumns.push({ index, dateStr });
    }
  });
  
  // 处理每一行（每行是一个关键词）
  for (let i = 1; i < rawData.length; i++) {
    const row = rawData[i];
    const keyword = row[0]; // 第一列是关键词
    
    if (!keyword || typeof keyword !== 'string') continue;
    
    const monthlyData = [];
    
    searchColumns.forEach(({ index, dateStr }) => {
      const volume = parseInt(row[index]) || 0;
      
      // 转换日期格式 "Jan 2022" -> "2022-01"
      const date = convertDateFormat(dateStr);
      if (date) {
        monthlyData.push({ date, volume });
      }
    });
    
    // 按日期排序
    monthlyData.sort((a, b) => a.date.localeCompare(b.date));
    
    keywords.push({
      keyword,
      source: 'google',
      monthlyVolume: monthlyData
    });
  }
  
  return keywords;
}

// 转换日期格式 "Jan 2022" -> "2022-01"
function convertDateFormat(dateStr) {
  const months = {
    'Jan': '01', 'Feb': '02', 'Mar': '03', 'Apr': '04',
    'May': '05', 'Jun': '06', 'Jul': '07', 'Aug': '08',
    'Sep': '09', 'Oct': '10', 'Nov': '11', 'Dec': '12'
  };
  
  const parts = dateStr.split(' ');
  if (parts.length !== 2) return null;
  
  const month = months[parts[0]];
  const year = parts[1];
  
  if (!month || !year) return null;
  
  return `${year}-${month}`;
}

// Parse KeywordHistory.xlsx (Amazon Keywords)
function parseAmazonKeywords(marketFolder) {
  const filePath = path.join(__dirname, '../data', marketFolder, 'KeywordHistory.xlsx');
  const workbook = XLSX.readFile(filePath);
  const keywords = [];
  
  workbook.SheetNames.forEach(sheetName => {
    const worksheet = workbook.Sheets[sheetName];
    const data = XLSX.utils.sheet_to_json(worksheet);
    
    const monthlyData = [];
    data.forEach(row => {
      const date = row.Date || row.date || row.Month || row.month || row['日期'];
      let volume = row['Search Volume'] || row.searchVolume || row.Volume || row.volume || row['搜索量'] || 0;
      
      // 确保volume是数字
      volume = parseInt(volume) || 0;
      
      if (date) {
        // 尝试解析日期
        let formattedDate;
        if (typeof date === 'string') {
          // 如果已经是 YYYY-MM 格式
          if (/^\d{4}-\d{2}/.test(date)) {
            formattedDate = date.substring(0, 7);
          } else {
            const dateObj = new Date(date);
            if (!isNaN(dateObj.getTime())) {
              const year = dateObj.getFullYear();
              const month = String(dateObj.getMonth() + 1).padStart(2, '0');
              formattedDate = `${year}-${month}`;
            }
          }
        } else if (typeof date === 'number') {
          // Excel日期序列号
          const excelDate = new Date((date - 25569) * 86400 * 1000);
          const year = excelDate.getFullYear();
          const month = String(excelDate.getMonth() + 1).padStart(2, '0');
          formattedDate = `${year}-${month}`;
        }
        
        if (formattedDate) {
          monthlyData.push({ date: formattedDate, volume });
        }
      }
    });
    
    // 按日期排序
    monthlyData.sort((a, b) => a.date.localeCompare(b.date));
    
    keywords.push({
      keyword: sheetName,
      source: 'amazon',
      monthlyVolume: monthlyData
    });
  });
  
  return keywords;
}

// Parse Brand-Social.xlsx
function parseBrandSocial(marketFolder) {
  const filePath = path.join(__dirname, '../data', marketFolder, 'Brand-Social.xlsx');
  const workbook = XLSX.readFile(filePath);
  const sheetName = workbook.SheetNames[0];
  const worksheet = workbook.Sheets[sheetName];
  const data = XLSX.utils.sheet_to_json(worksheet);
  
  return data.map(row => ({
    brand: row.Brand || row['品牌'] || row.brand || 'Unknown',
    website: row.Website || row['独立站'] || row.website || '',
    youtube: {
      url: row['YouTube URL'] || row['YouTube链接'] || row.youtubeUrl || row.Youtube || '',
      subscribers: parseInt(row['YouTube Subscribers'] || row['YouTube订阅数'] || row.youtubeSubscribers || 0)
    },
    instagram: {
      url: row['Instagram URL'] || row['Instagram链接'] || row.instagramUrl || row.Instagram || '',
      followers: parseInt(row['Instagram Followers'] || row['Instagram粉丝数'] || row.instagramFollowers || 0)
    },
    facebook: {
      url: row['Facebook URL'] || row['Facebook链接'] || row.facebookUrl || row.Facebook || '',
      followers: parseInt(row['Facebook Followers/Likes'] || row['Facebook Followers'] || row['Facebook粉丝数'] || row.facebookFollowers || 0)
    },
    reddit: {
      url: row['Reddit URL/Search'] || row['Reddit URL'] || row['Reddit链接'] || row.redditUrl || row.Reddit || '',
      posts: parseInt(row['Reddit Mentions (approx)'] || row['Reddit Posts'] || row['Reddit讨论数'] || row.redditPosts || 0)
    }
  }));
}

// Calculate market metrics
function calculateMarketMetrics(products, productSales, keywords) {
  // Calculate revenue from sales data
  const monthlySales = {};
  
  Object.values(productSales).forEach(sales => {
    if (sales.monthlySales) {
      sales.monthlySales.forEach(month => {
        if (!monthlySales[month.date]) {
          monthlySales[month.date] = 0;
        }
        monthlySales[month.date] += month.sales || 0;
      });
    }
  });
  
  const sortedMonths = Object.keys(monthlySales).sort();
  
  // Calculate total revenue (last 12 months)
  let totalRevenue = 0;
  if (sortedMonths.length >= 12) {
    const last12Months = sortedMonths.slice(-12);
    totalRevenue = last12Months.reduce((sum, month) => sum + monthlySales[month], 0);
  } else if (sortedMonths.length > 0) {
    totalRevenue = sortedMonths.reduce((sum, month) => sum + monthlySales[month], 0);
  }
  
  // Calculate CAGR (Compound Annual Growth Rate)
  let cagr = null; // null表示无法计算
  if (sortedMonths.length >= 12) {
    const firstYearRevenue = sortedMonths.slice(0, 12).reduce((sum, month) => sum + monthlySales[month], 0);
    const lastYearRevenue = sortedMonths.slice(-12).reduce((sum, month) => sum + monthlySales[month], 0);
    
    // 只有当两个时期都有实际销售数据时才计算CAGR
    // 同时设置最小阈值以避免基数太小导致的极端值
    if (firstYearRevenue > 1000 && lastYearRevenue > 0) {
      // Calculate number of years between first and last 12-month periods
      // Using the midpoint of each period for more accurate calculation
      const firstDate = new Date(sortedMonths[5]); // Midpoint of first 12 months
      const lastDate = new Date(sortedMonths[sortedMonths.length - 6]); // Midpoint of last 12 months
      const years = (lastDate - firstDate) / (365.25 * 24 * 60 * 60 * 1000);
      
      if (years > 0.5) { // 至少需要半年的时间跨度
        // CAGR = (Ending Value / Beginning Value)^(1/years) - 1
        cagr = (Math.pow(lastYearRevenue / firstYearRevenue, 1 / years) - 1) * 100;
        
        // 对极端值进行合理限制（-99%到+999%）
        if (cagr < -99) cagr = -99;
        if (cagr > 999) cagr = 999;
      }
    }
  }
  
  // Calculate monthly trends
  const monthlyTrends = sortedMonths.map(date => ({
    date,
    revenue: monthlySales[date] || 0
  }));
  
  // Calculate keyword trends
  const keywordTrends = {};
  keywords.forEach(keyword => {
    keyword.monthlyVolume.forEach(month => {
      if (!keywordTrends[month.date]) {
        keywordTrends[month.date] = 0;
      }
      keywordTrends[month.date] += month.volume || 0;
    });
  });
  
  const sortedKeywordMonths = Object.keys(keywordTrends).sort();
  const keywordMonthlyTrends = sortedKeywordMonths.map(date => ({
    date,
    volume: keywordTrends[date] || 0
  }));
  
  // Calculate total search volume (latest month)
  let totalSearchVolume = 0;
  if (sortedKeywordMonths.length > 0) {
    const latestMonth = sortedKeywordMonths[sortedKeywordMonths.length - 1];
    totalSearchVolume = keywordTrends[latestMonth] || 0;
  }
  
  // Count brands
  const brands = new Set();
  products.forEach(product => {
    if (product.brand && product.brand !== 'Unknown') {
      brands.add(product.brand);
    }
  });
  
  return {
    totalRevenue,
    totalProducts: products.length,
    brandCount: brands.size,
    searchVolume: totalSearchVolume,
    cagr,
    monthlyTrends,
    keywordMonthlyTrends
  };
}

// Calculate brand metrics
function calculateBrandMetrics(brand, products, productSales) {
  const brandProducts = products.filter(p => p.brand === brand);
  
  const monthlySales = {};
  
  brandProducts.forEach(product => {
    const sales = productSales[product.asin];
    if (sales && sales.monthlySales) {
      sales.monthlySales.forEach(month => {
        if (!monthlySales[month.date]) {
          monthlySales[month.date] = 0;
        }
        monthlySales[month.date] += month.sales || 0;
      });
    }
  });
  
  const sortedMonths = Object.keys(monthlySales).sort();
  
  // Calculate total revenue (last 12 months)
  let totalRevenue = 0;
  if (sortedMonths.length >= 12) {
    const last12Months = sortedMonths.slice(-12);
    totalRevenue = last12Months.reduce((sum, month) => sum + monthlySales[month], 0);
  } else if (sortedMonths.length > 0) {
    totalRevenue = sortedMonths.reduce((sum, month) => sum + monthlySales[month], 0);
  }
  
  // Calculate CAGR (Compound Annual Growth Rate)
  let cagr = null; // null表示无法计算
  if (sortedMonths.length >= 12) {
    // Use at least 12 months for each period
    const firstYearRevenue = sortedMonths.slice(0, 12).reduce((sum, month) => sum + monthlySales[month], 0);
    const lastYearRevenue = sortedMonths.slice(-12).reduce((sum, month) => sum + monthlySales[month], 0);
    
    // 只有当两个时期都有实际销售数据时才计算CAGR
    // 同时设置最小阈值（例如$1000）以避免基数太小导致的极端值
    if (firstYearRevenue > 1000 && lastYearRevenue > 0) {
      // Calculate number of years between first and last 12-month periods
      const firstDate = new Date(sortedMonths[5]); // Midpoint of first 12 months
      const lastDate = new Date(sortedMonths[sortedMonths.length - 6]); // Midpoint of last 12 months
      const years = (lastDate - firstDate) / (365.25 * 24 * 60 * 60 * 1000);
      
      if (years > 0.5) { // 至少需要半年的时间跨度
        // CAGR = (Ending Value / Beginning Value)^(1/years) - 1
        cagr = (Math.pow(lastYearRevenue / firstYearRevenue, 1 / years) - 1) * 100;
        
        // 对极端值进行合理限制（-99%到+999%）
        if (cagr < -99) cagr = -99;
        if (cagr > 999) cagr = 999;
      }
    }
  }
  
  const monthlyTrends = sortedMonths.map(date => ({
    date,
    revenue: monthlySales[date] || 0
  }));
  
  return {
    brand,
    totalRevenue,
    productCount: brandProducts.length,
    cagr,
    monthlyTrends
  };
}

// Main processing function
function processMarket(market) {
  console.log(`Processing ${market.name}...`);
  
  const products = parseProducts(market.folder);
  const productSales = parseProductSales(market.folder);
  const googleKeywords = parseGoogleKeywords(market.folder);
  const amazonKeywords = parseAmazonKeywords(market.folder);
  const brandSocial = parseBrandSocial(market.folder);
  
  const allKeywords = [...googleKeywords, ...amazonKeywords];
  
  const metrics = calculateMarketMetrics(products, productSales, allKeywords);
  
  console.log(`  - Products: ${products.length}`);
  console.log(`  - Monthly trends: ${metrics.monthlyTrends.length} months`);
  console.log(`  - Keyword trends: ${metrics.keywordMonthlyTrends.length} months`);
  console.log(`  - Total revenue: $${metrics.totalRevenue.toLocaleString()}`);
  
  // Calculate brand metrics
  const brands = [...new Set(products.map(p => p.brand).filter(b => b && b !== 'Unknown'))];
  const brandMetrics = brands.map(brand => {
    const metrics = calculateBrandMetrics(brand, products, productSales);
    const social = brandSocial.find(b => b.brand === brand) || {};
    return {
      ...metrics,
      ...social
    };
  }).sort((a, b) => b.totalRevenue - a.totalRevenue);
  
  return {
    id: market.id,
    name: market.name,
    metrics,
    products,
    productSales,
    keywords: allKeywords,
    brands: brandMetrics
  };
}

// Process all markets and save to JSON
function main() {
  const outputDir = path.join(__dirname, '../public/data');
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }
  
  const allMarkets = markets.map(market => {
    const marketData = processMarket(market);
    
    // Save individual market data
    fs.writeFileSync(
      path.join(outputDir, `${market.id}.json`),
      JSON.stringify(marketData, null, 2)
    );
    
    return {
      id: marketData.id,
      name: marketData.name,
      metrics: marketData.metrics
    };
  });
  
  // Save markets summary
  fs.writeFileSync(
    path.join(outputDir, 'markets.json'),
    JSON.stringify(allMarkets, null, 2)
  );
  
  console.log('\n✅ Data processing complete!');
}

main();
