module.exports = {
  presets: ['module:@react-native/babel-preset'],
  plugins: [
    [
      'module-resolver',
      {
        root: ['./src'],
        alias: {
          '@': './src',
          '@components': './src/components',
          '@screens': './src/screens',
          '@services': './src/services',
          '@utils': './src/utils',
          '@types': './src/types',
          '@constants': './src/constants',
          '@hooks': './src/hooks',
          '@store': './src/store',
          '@navigation': './src/navigation',
        },
      },
    ],
  ],
};
