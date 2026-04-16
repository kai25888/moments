import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.moments.plus',
  appName: 'Moments',
  webDir: '.output/public',
  server: {
    url: 'https://moments.979569933.xyz',
    androidScheme: 'https',
    cleartext: true,
  },
  plugins: {
    SplashScreen: {
      launchShowDuration: 1000,
      launchAutoHide: true,
      backgroundColor: '#18181b',
      showSpinner: false,
    },
    StatusBar: {
      style: 'DARK',
      backgroundColor: '#18181b',
    },
    Keyboard: {
      resize: 'none',
    },
  },
};

export default config;
