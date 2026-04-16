/**
 * Capacitor 原生桥接插件
 * 处理安全区域、状态栏、Android 返回键等原生能力
 */
import { App } from '@capacitor/app';
import { StatusBar, Style } from '@capacitor/status-bar';
import { SplashScreen } from '@capacitor/splash-screen';
import { Capacitor } from '@capacitor/core';

export function useNativeBridge() {
  const isNative = ref(false);

  onMounted(async () => {
    isNative.value = Capacitor.isNativePlatform();

    if (!isNative.value) return;

    // 1. 配置状态栏
    try {
      await StatusBar.setStyle({ style: Style.Light });
      await StatusBar.setBackgroundColor({ color: '#18181b' });
      await StatusBar.setOverlaysWebView({ overlay: true });
    } catch (e) {
      console.warn('[Native] StatusBar config failed:', e);
    }

    // 2. 隐藏启动屏
    try {
      await SplashScreen.hide({ fadeOutDuration: 300 });
    } catch (e) {
      console.warn('[Native] SplashScreen hide failed:', e);
    }

    // 3. Android 返回键处理
    if (Capacitor.getPlatform() === 'android') {
      App.addListener('backButton', ({ canGoBack }) => {
        if (canGoBack) {
          window.history.back();
        } else {
          App.exitApp();
        }
      });
    }
  });

  return { isNative };
}
