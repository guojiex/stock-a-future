/**
 * 基本面数据展示组件
 */

import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, ActivityIndicator } from 'react-native';
import { Card, useTheme, Button, SegmentedButtons, Divider } from 'react-native-paper';
import { apiService } from '@/services/apiService';
import { SkeletonCard } from '@/components/common/SkeletonLoader';

interface FundamentalDataProps {
  stockCode: string;
  stockName: string;
}

interface StockBasicInfo {
  ts_code: string;
  name: string;
  area?: string;
  industry?: string;
  fullname?: string;
  enname?: string;
  market?: string;
  exchange?: string;
  list_date?: string;
  delist_date?: string;
}

interface DailyBasicInfo {
  close?: number;
  pe?: number;
  pb?: number;
  ps?: number;
  dv_ratio?: number;
  total_share?: number;
  float_share?: number;
  free_share?: number;
  total_mv?: number;
  circ_mv?: number;
}

interface FinancialStatement {
  end_date: string;
  report_type?: string;
  [key: string]: any;
}

interface FundamentalDataResponse {
  stock_basic?: StockBasicInfo;
  daily_basic?: DailyBasicInfo;
  income_statement?: FinancialStatement;
  balance_sheet?: FinancialStatement;
  cash_flow_statement?: FinancialStatement;
}

const FundamentalData: React.FC<FundamentalDataProps> = ({ stockCode, stockName }) => {
  const theme = useTheme();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<FundamentalDataResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [selectedPeriod, setSelectedPeriod] = useState('annual'); // annual, quarter

  // 加载基本面数据
  const loadFundamentalData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiService.getFundamentalData(stockCode, selectedPeriod);
      
      if (response.success && response.data) {
        setData(response.data);
      } else {
        setError(response.error || '加载基本面数据失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载基本面数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFundamentalData();
  }, [stockCode, selectedPeriod]);

  // 格式化数值
  const formatNumber = (value: any, unit?: string) => {
    if (value === null || value === undefined || value === '') return 'N/A';
    
    const num = parseFloat(value);
    if (isNaN(num)) return 'N/A';
    
    if (unit === 'percent') {
      return `${num.toFixed(2)}%`;
    } else if (unit === 'wan') {
      return `${(num / 10000).toFixed(2)}万`;
    } else if (unit === 'yi') {
      return `${(num / 100000000).toFixed(2)}亿`;
    } else {
      return num.toLocaleString();
    }
  };

  // 格式化日期
  const formatDate = (dateString: string) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('zh-CN');
    } catch {
      return dateString;
    }
  };

  // 渲染信息网格
  const renderInfoGrid = (title: string, items: Array<{ label: string; value: any; unit?: string }>) => (
    <Card style={styles.sectionCard} key={title}>
      <Card.Content>
        <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
          {title}
        </Text>
        <View style={styles.infoGrid}>
          {items.map((item, index) => (
            <View key={index} style={styles.infoItem}>
              <Text style={[styles.infoLabel, { color: theme.colors.onSurfaceVariant }]}>
                {item.label}
              </Text>
              <Text style={[styles.infoValue, { color: theme.colors.onSurface }]}>
                {formatNumber(item.value, item.unit)}
              </Text>
            </View>
          ))}
        </View>
      </Card.Content>
    </Card>
  );

  // 渲染财务报表
  const renderFinancialTable = (title: string, statement: FinancialStatement | undefined, fields: Array<{ key: string; label: string; unit?: string }>) => {
    if (!statement) return null;

    return (
      <Card style={styles.sectionCard} key={title}>
        <Card.Content>
          <Text style={[styles.sectionTitle, { color: theme.colors.onSurface }]}>
            {title}
          </Text>
          <Text style={[styles.reportPeriod, { color: theme.colors.onSurfaceVariant }]}>
            报告期: {formatDate(statement.end_date)}
          </Text>
          <Divider style={styles.divider} />
          
          <View style={styles.tableContainer}>
            {fields.map((field, index) => (
              <View key={index} style={styles.tableRow}>
                <Text style={[styles.tableLabel, { color: theme.colors.onSurfaceVariant }]}>
                  {field.label}
                </Text>
                <Text style={[styles.tableValue, { color: theme.colors.onSurface }]}>
                  {formatNumber(statement[field.key], field.unit)}
                </Text>
              </View>
            ))}
          </View>
        </Card.Content>
      </Card>
    );
  };

  if (loading) {
    return (
      <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: theme.colors.onSurface }]}>
            基本面数据
          </Text>
          <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
            {stockName} ({stockCode})
          </Text>
        </View>
        <View style={styles.periodContainer}>
          <SegmentedButtons
            value={selectedPeriod}
            onValueChange={setSelectedPeriod}
            buttons={[
              { value: 'annual', label: '年报' },
              { value: 'quarter', label: '季报' },
            ]}
            style={styles.segmentedButtons}
          />
        </View>
        <View style={styles.contentContainer}>
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
        </View>
      </ScrollView>
    );
  }

  if (error) {
    return (
      <View style={styles.centerContainer}>
        <Text style={[styles.errorText, { color: theme.colors.error }]}>
          {error}
        </Text>
        <Button mode="outlined" onPress={loadFundamentalData} style={styles.retryButton}>
          重试
        </Button>
      </View>
    );
  }

  if (!data) {
    return (
      <View style={styles.centerContainer}>
        <Text style={[styles.emptyText, { color: theme.colors.onSurfaceVariant }]}>
          暂无基本面数据
        </Text>
        <Button mode="outlined" onPress={loadFundamentalData} style={styles.retryButton}>
          刷新
        </Button>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container} showsVerticalScrollIndicator={false}>
      <View style={styles.header}>
        <Text style={[styles.title, { color: theme.colors.onSurface }]}>
          基本面数据
        </Text>
        <Text style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}>
          {stockName} ({stockCode})
        </Text>
      </View>

      {/* 周期选择 */}
      <View style={styles.periodContainer}>
        <SegmentedButtons
          value={selectedPeriod}
          onValueChange={setSelectedPeriod}
          buttons={[
            { value: 'annual', label: '年报' },
            { value: 'quarter', label: '季报' },
          ]}
          style={styles.segmentedButtons}
        />
      </View>

      <View style={styles.contentContainer}>
        {/* 股票基本信息 */}
        {data.stock_basic && renderInfoGrid(
          '📋 股票基本信息',
          [
            { label: '股票名称', value: data.stock_basic.name },
            { label: '所属地域', value: data.stock_basic.area },
            { label: '所属行业', value: data.stock_basic.industry },
            { label: '上市日期', value: formatDate(data.stock_basic.list_date || '') },
            { label: '交易所', value: data.stock_basic.exchange },
            { label: '市场类型', value: data.stock_basic.market },
          ]
        )}

        {/* 每日基本面指标 */}
        {data.daily_basic && renderInfoGrid(
          '📊 每日基本面指标',
          [
            { label: '市盈率', value: data.daily_basic.pe, unit: 'percent' },
            { label: '市净率', value: data.daily_basic.pb, unit: 'percent' },
            { label: '市销率', value: data.daily_basic.ps, unit: 'percent' },
            { label: '股息率', value: data.daily_basic.dv_ratio, unit: 'percent' },
            { label: '总市值', value: data.daily_basic.total_mv, unit: 'wan' },
            { label: '流通市值', value: data.daily_basic.circ_mv, unit: 'wan' },
          ]
        )}

        {/* 利润表 */}
        {renderFinancialTable(
          '💰 利润表',
          data.income_statement,
          [
            { key: 'total_revenue', label: '营业总收入', unit: 'wan' },
            { key: 'revenue', label: '营业收入', unit: 'wan' },
            { key: 'oper_cost', label: '营业成本', unit: 'wan' },
            { key: 'operate_profit', label: '营业利润', unit: 'wan' },
            { key: 'total_profit', label: '利润总额', unit: 'wan' },
            { key: 'n_income', label: '净利润', unit: 'wan' },
            { key: 'basic_eps', label: '基本每股收益' },
            { key: 'diluted_eps', label: '稀释每股收益' },
          ]
        )}

        {/* 资产负债表 */}
        {renderFinancialTable(
          '📋 资产负债表',
          data.balance_sheet,
          [
            { key: 'total_assets', label: '资产总计', unit: 'wan' },
            { key: 'total_cur_assets', label: '流动资产合计', unit: 'wan' },
            { key: 'monetary_cap', label: '货币资金', unit: 'wan' },
            { key: 'accounts_receiv', label: '应收账款', unit: 'wan' },
            { key: 'total_liab', label: '负债合计', unit: 'wan' },
            { key: 'total_cur_liab', label: '流动负债合计', unit: 'wan' },
            { key: 'total_hldr_eqy_exc_min_int', label: '所有者权益合计', unit: 'wan' },
            { key: 'cap_rese', label: '资本公积', unit: 'wan' },
          ]
        )}

        {/* 现金流量表 */}
        {renderFinancialTable(
          '💵 现金流量表',
          data.cash_flow_statement,
          [
            { key: 'n_cashflow_act', label: '经营活动现金流净额', unit: 'wan' },
            { key: 'cash_recp_sg_and_rs', label: '销售商品收到现金', unit: 'wan' },
            { key: 'cash_pay_goods_purch', label: '购买商品支付现金', unit: 'wan' },
            { key: 'n_cashflow_inv_act', label: '投资活动现金流净额', unit: 'wan' },
            { key: 'n_cashflow_fin_act', label: '筹资活动现金流净额', unit: 'wan' },
            { key: 'n_cash_flows_fnc_act', label: '筹资活动现金流净额', unit: 'wan' },
            { key: 'cash_at_eop', label: '期末现金余额', unit: 'wan' },
          ]
        )}
      </View>

      <View style={styles.refreshContainer}>
        <Button mode="outlined" onPress={loadFundamentalData} disabled={loading}>
          刷新数据
        </Button>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  header: {
    padding: 16,
    paddingBottom: 8,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 14,
  },
  periodContainer: {
    paddingHorizontal: 16,
    marginBottom: 16,
  },
  segmentedButtons: {
    width: '100%',
  },
  contentContainer: {
    padding: 16,
    paddingTop: 0,
    gap: 16,
  },
  sectionCard: {
    marginBottom: 8,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 12,
  },
  reportPeriod: {
    fontSize: 12,
    marginBottom: 8,
  },
  divider: {
    marginVertical: 8,
  },
  infoGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 12,
  },
  infoItem: {
    minWidth: '45%',
    marginBottom: 12,
  },
  infoLabel: {
    fontSize: 12,
    marginBottom: 4,
  },
  infoValue: {
    fontSize: 14,
    fontWeight: '600',
  },
  tableContainer: {
    gap: 8,
  },
  tableRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#e0e0e0',
  },
  tableLabel: {
    fontSize: 13,
    flex: 1,
  },
  tableValue: {
    fontSize: 13,
    fontWeight: '600',
    textAlign: 'right',
    minWidth: 80,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 14,
  },
  errorText: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 16,
  },
  emptyText: {
    fontSize: 14,
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    marginTop: 8,
  },
  refreshContainer: {
    padding: 16,
    alignItems: 'center',
  },
});

export default FundamentalData;
