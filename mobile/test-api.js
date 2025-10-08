// 测试React Native API调用 - 陕西煤业(601225)数据调试
const fetch = require('node-fetch');

async function testShanxiCoalAPI() {
  const baseURL = 'http://127.0.0.1:8081/api/v1/';
  const stockCode = '601225'; // 陕西煤业
  
  // 测试多个时间范围
  const testCases = [
    { days: 90, label: '3个月' },
    { days: 180, label: '半年' },
    { days: 365, label: '1年' }
  ];
  
  for (const testCase of testCases) {
    console.log('\n' + '='.repeat(80));
    console.log(`🔍 测试陕西煤业(${stockCode}) - 时间范围: ${testCase.label} (${testCase.days}天)`);
    console.log('='.repeat(80));
    
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - testCase.days);
    
    const formatDate = (date) => {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}${month}${day}`;
    };
    
    const start = formatDate(startDate);
    const end = formatDate(endDate);
    const adjust = 'qfq'; // 前复权
    
    const url = `${baseURL}stocks/${stockCode}/daily?start_date=${start}&end_date=${end}&adjust=${adjust}`;
    
    console.log('📡 API请求:', {
      url,
      start_date: start,
      end_date: end,
      adjust
    });
    
    try {
      const response = await fetch(url);
      const data = await response.json();
      
      console.log('📥 API响应状态:', response.status);
      console.log('✅ 响应成功:', data.success);
      
      if (data.success && data.data && Array.isArray(data.data)) {
        const dataArray = data.data;
        
        console.log('\n📊 数据统计:');
        console.log('  - 总记录数:', dataArray.length);
        console.log('  - 日期范围:', dataArray[0]?.trade_date, '至', dataArray[dataArray.length - 1]?.trade_date);
        
        // 计算价格统计
        const highs = dataArray.map(d => parseFloat(d.high));
        const lows = dataArray.map(d => parseFloat(d.low));
        const closes = dataArray.map(d => parseFloat(d.close));
        
        const maxHigh = Math.max(...highs);
        const minLow = Math.min(...lows);
        const firstClose = closes[0];
        const lastClose = closes[closes.length - 1];
        
        console.log('\n📈 价格统计:');
        console.log('  - 最高价:', maxHigh.toFixed(2), '元');
        console.log('  - 最低价:', minLow.toFixed(2), '元');
        console.log('  - 期初收盘:', firstClose.toFixed(2), '元');
        console.log('  - 期末收盘:', lastClose.toFixed(2), '元');
        console.log('  - 价格区间:', `${minLow.toFixed(2)} - ${maxHigh.toFixed(2)} 元`);
        
        // 找到最高价对应的日期
        const maxHighIndex = highs.indexOf(maxHigh);
        const maxHighDate = dataArray[maxHighIndex];
        console.log('  - 最高价日期:', maxHighDate.trade_date);
        console.log('  - 最高价详情:', {
          日期: maxHighDate.trade_date,
          开盘: parseFloat(maxHighDate.open).toFixed(2),
          最高: parseFloat(maxHighDate.high).toFixed(2),
          最低: parseFloat(maxHighDate.low).toFixed(2),
          收盘: parseFloat(maxHighDate.close).toFixed(2)
        });
        
        // 找到最低价对应的日期
        const minLowIndex = lows.indexOf(minLow);
        const minLowDate = dataArray[minLowIndex];
        console.log('  - 最低价日期:', minLowDate.trade_date);
        console.log('  - 最低价详情:', {
          日期: minLowDate.trade_date,
          开盘: parseFloat(minLowDate.open).toFixed(2),
          最高: parseFloat(minLowDate.high).toFixed(2),
          最低: parseFloat(minLowDate.low).toFixed(2),
          收盘: parseFloat(minLowDate.close).toFixed(2)
        });
        
        // 显示前5条和后5条数据样本
        console.log('\n📄 数据样本 (前5条):');
        dataArray.slice(0, 5).forEach((item, index) => {
          console.log(`  ${index + 1}. ${item.trade_date} - 开:${parseFloat(item.open).toFixed(2)} 高:${parseFloat(item.high).toFixed(2)} 低:${parseFloat(item.low).toFixed(2)} 收:${parseFloat(item.close).toFixed(2)}`);
        });
        
        console.log('\n📄 数据样本 (后5条):');
        dataArray.slice(-5).forEach((item, index) => {
          console.log(`  ${dataArray.length - 4 + index}. ${item.trade_date} - 开:${parseFloat(item.open).toFixed(2)} 高:${parseFloat(item.high).toFixed(2)} 低:${parseFloat(item.low).toFixed(2)} 收:${parseFloat(item.close).toFixed(2)}`);
        });
        
        // 检查是否有价格异常（超过17元）
        const highPrices = dataArray.filter(d => parseFloat(d.high) >= 17);
        if (highPrices.length > 0) {
          console.log('\n⚠️  警告: 发现价格 >= 17元的数据:');
          highPrices.forEach(item => {
            console.log(`  - ${item.trade_date}: 最高价 ${parseFloat(item.high).toFixed(2)}元`, item);
          });
        } else {
          console.log('\n✅ 未发现价格 >= 17元的数据');
        }
        
      } else {
        console.log('❌ 响应数据无效:', data);
      }
      
    } catch (error) {
      console.error('❌ 请求失败:', error.message);
    }
  }
}

console.log('🚀 开始测试陕西煤业API数据...\n');
testShanxiCoalAPI().then(() => {
  console.log('\n✅ 测试完成');
}).catch(error => {
  console.error('\n❌ 测试失败:', error);
});
