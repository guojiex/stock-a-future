/**
 * Stock-A-Future React Native App
 * 主入口文件
 */

import React from 'react';
import { StatusBar, useColorScheme } from 'react-native';
import { Provider } from 'react-redux';
import { PaperProvider } from 'react-native-paper';
import { SafeAreaProvider } from 'react-native-safe-area-context';

import AppNavigator from '@/navigation/AppNavigator';
import store from '@/store';
import { lightTheme, darkTheme } from '@/constants/themes';

function App(): JSX.Element {
  const isDarkMode = useColorScheme() === 'dark';

  // 根据系统主题选择应用主题
  const theme = isDarkMode ? darkTheme : lightTheme;

  return (
    <Provider store={store}>
      <PaperProvider theme={theme}>
        <SafeAreaProvider>
          <StatusBar
            barStyle={isDarkMode ? 'light-content' : 'dark-content'}
            backgroundColor={theme.colors.surface}
          />
          <AppNavigator />
        </SafeAreaProvider>
      </PaperProvider>
    </Provider>
  );
}

export default App;
