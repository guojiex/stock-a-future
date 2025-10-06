/**
 * 应用主导航器
 */

import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import Icon from 'react-native-vector-icons/MaterialIcons';

// 导入页面组件
import MarketScreen from '@/screens/Market/MarketScreen';
import StockDetailScreen from '@/screens/Market/StockDetailScreen';
import TechnicalIndicatorsScreen from '@/screens/Market/TechnicalIndicatorsScreen';
import FundamentalScreen from '@/screens/Market/FundamentalScreen';

import SearchScreen from '@/screens/Search/SearchScreen';
import SearchResultScreen from '@/screens/Search/SearchResultScreen';

import FavoritesScreen from '@/screens/Favorites/FavoritesScreen';
import FavoritesManageScreen from '@/screens/Favorites/FavoritesManageScreen';
import GroupManageScreen from '@/screens/Favorites/GroupManageScreen';

import BacktestScreen from '@/screens/Backtest/BacktestScreen';
import BacktestResultScreen from '@/screens/Backtest/BacktestResultScreen';
import StrategyEditScreen from '@/screens/Backtest/StrategyEditScreen';

import SettingsScreen from '@/screens/Settings/SettingsScreen';

// 导航参数类型定义
export type RootTabParamList = {
  MarketTab: undefined;
  SearchTab: undefined;
  FavoritesTab: undefined;
  BacktestTab: undefined;
  SettingsTab: undefined;
};

export type MarketStackParamList = {
  MarketHome: undefined;
  StockDetail: {
    stockCode: string;
    stockName: string;
  };
  TechnicalIndicators: {
    stockCode: string;
    stockName: string;
  };
  Fundamental: {
    stockCode: string;
    stockName: string;
  };
};

export type SearchStackParamList = {
  SearchHome: undefined;
  SearchResult: {
    query: string;
  };
};

export type FavoritesStackParamList = {
  FavoritesHome: undefined;
  FavoritesManage: undefined;
  GroupManage: {
    groupId?: number;
  };
};

export type BacktestStackParamList = {
  BacktestHome: undefined;
  BacktestResult: {
    backtestId: string;
  };
  StrategyEdit: {
    strategyId?: string;
  };
};

export type SettingsStackParamList = {
  SettingsHome: undefined;
};

// 创建导航器
const Tab = createBottomTabNavigator<RootTabParamList>();
const MarketStack = createNativeStackNavigator<MarketStackParamList>();
const SearchStack = createNativeStackNavigator<SearchStackParamList>();
const FavoritesStack = createNativeStackNavigator<FavoritesStackParamList>();
const BacktestStack = createNativeStackNavigator<BacktestStackParamList>();
const SettingsStack = createNativeStackNavigator<SettingsStackParamList>();

// 市场页面栈导航
function MarketStackNavigator() {
  return (
    <MarketStack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      <MarketStack.Screen
        name="MarketHome"
        component={MarketScreen}
        options={{
          title: '市场',
          headerShown: false, // 在Tab中显示标题
        }}
      />
      <MarketStack.Screen
        name="StockDetail"
        component={StockDetailScreen}
        options={({ route }) => ({
          title: route.params.stockName || route.params.stockCode,
          headerBackTitle: '返回',
        })}
      />
      <MarketStack.Screen
        name="TechnicalIndicators"
        component={TechnicalIndicatorsScreen}
        options={({ route }) => ({
          title: `${route.params.stockName || route.params.stockCode} - 技术指标`,
          headerBackTitle: '返回',
        })}
      />
      <MarketStack.Screen
        name="Fundamental"
        component={FundamentalScreen}
        options={({ route }) => ({
          title: `${route.params.stockName || route.params.stockCode} - 基本面`,
          headerBackTitle: '返回',
        })}
      />
    </MarketStack.Navigator>
  );
}

// 搜索页面栈导航
function SearchStackNavigator() {
  return (
    <SearchStack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      <SearchStack.Screen
        name="SearchHome"
        component={SearchScreen}
        options={{
          title: '搜索',
          headerShown: false, // 在Tab中显示标题
        }}
      />
      <SearchStack.Screen
        name="SearchResult"
        component={SearchResultScreen}
        options={({ route }) => ({
          title: `搜索: ${route.params.query}`,
          headerBackTitle: '返回',
        })}
      />
    </SearchStack.Navigator>
  );
}

// 收藏页面栈导航
function FavoritesStackNavigator() {
  return (
    <FavoritesStack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      <FavoritesStack.Screen
        name="FavoritesHome"
        component={FavoritesScreen}
        options={{
          title: '收藏',
          headerShown: false, // 在Tab中显示标题
        }}
      />
      <FavoritesStack.Screen
        name="FavoritesManage"
        component={FavoritesManageScreen}
        options={{
          title: '管理收藏',
          headerBackTitle: '返回',
        }}
      />
      <FavoritesStack.Screen
        name="GroupManage"
        component={GroupManageScreen}
        options={({ route }) => ({
          title: route.params.groupId ? '编辑分组' : '新建分组',
          headerBackTitle: '返回',
        })}
      />
    </FavoritesStack.Navigator>
  );
}

// 回测页面栈导航
function BacktestStackNavigator() {
  return (
    <BacktestStack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      <BacktestStack.Screen
        name="BacktestHome"
        component={BacktestScreen}
        options={{
          title: '回测',
          headerShown: false, // 在Tab中显示标题
        }}
      />
      <BacktestStack.Screen
        name="BacktestResult"
        component={BacktestResultScreen}
        options={{
          title: '回测结果',
          headerBackTitle: '返回',
        }}
      />
      <BacktestStack.Screen
        name="StrategyEdit"
        component={StrategyEditScreen}
        options={({ route }) => ({
          title: route.params.strategyId ? '编辑策略' : '新建策略',
          headerBackTitle: '返回',
        })}
      />
    </BacktestStack.Navigator>
  );
}

// 设置页面栈导航
function SettingsStackNavigator() {
  return (
    <SettingsStack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      <SettingsStack.Screen
        name="SettingsHome"
        component={SettingsScreen}
        options={{
          title: '设置',
          headerShown: false, // 在Tab中显示标题
        }}
      />
    </SettingsStack.Navigator>
  );
}

// 主Tab导航器
function TabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: string;

          switch (route.name) {
            case 'MarketTab':
              iconName = focused ? 'trending-up' : 'trending-up';
              break;
            case 'SearchTab':
              iconName = focused ? 'search' : 'search';
              break;
            case 'FavoritesTab':
              iconName = focused ? 'star' : 'star-border';
              break;
            case 'BacktestTab':
              iconName = focused ? 'bar-chart' : 'bar-chart';
              break;
            case 'SettingsTab':
              iconName = focused ? 'settings' : 'settings';
              break;
            default:
              iconName = 'help';
              break;
          }

          return <Icon name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: '#007AFF',
        tabBarInactiveTintColor: '#8E8E93',
        tabBarStyle: {
          backgroundColor: '#F8F9FA',
          borderTopWidth: 1,
          borderTopColor: '#E1E1E1',
          paddingBottom: 5,
          paddingTop: 5,
          height: 60,
        },
        tabBarLabelStyle: {
          fontSize: 12,
          fontWeight: '500',
        },
        headerStyle: {
          backgroundColor: '#f8f9fa',
        },
        headerTintColor: '#333',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      })}
    >
      <Tab.Screen
        name="MarketTab"
        component={MarketStackNavigator}
        options={{
          title: '市场',
          tabBarLabel: '市场',
        }}
      />
      <Tab.Screen
        name="SearchTab"
        component={SearchStackNavigator}
        options={{
          title: '搜索',
          tabBarLabel: '搜索',
        }}
      />
      <Tab.Screen
        name="FavoritesTab"
        component={FavoritesStackNavigator}
        options={{
          title: '收藏',
          tabBarLabel: '收藏',
        }}
      />
      <Tab.Screen
        name="BacktestTab"
        component={BacktestStackNavigator}
        options={{
          title: '回测',
          tabBarLabel: '回测',
        }}
      />
      <Tab.Screen
        name="SettingsTab"
        component={SettingsStackNavigator}
        options={{
          title: '设置',
          tabBarLabel: '设置',
        }}
      />
    </Tab.Navigator>
  );
}

// 应用主导航器
export default function AppNavigator() {
  return (
    <NavigationContainer>
      <TabNavigator />
    </NavigationContainer>
  );
}
