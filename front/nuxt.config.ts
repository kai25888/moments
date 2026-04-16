// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: false},
    modules: [
        "@nuxt/ui", 
        '@nuxt/icon', 
        '@nuxtjs/color-mode', 
        '@vueuse/nuxt', 
        'dayjs-nuxt',
        '@nuxt/image'  // 🔥 图片优化模块
    ],
    colorMode: {
        preference: 'system',  // 自动跟随系统主题
        fallback: 'light',    // 系统不支持时的后备值
        classSuffix: '',
        dataValue: 'theme',
        storageKey: 'nuxt-color-mode',
    },
    ssr: false,
    dayjs: {
        locales: ['zh'],
        defaultLocale: 'zh'
    },
    icon: {
        clientBundle: {
            scan: {
                globInclude: ['**/*.{vue,jsx,tsx}', 'node_modules/@nuxt/ui/**/*.js'],
                globExclude: ['.*', 'coverage', 'test', 'tests', 'dist', 'build'],
            },
        },
    },
    tailwindcss: {
        safelist: [
            'grid-cols-1',
            'grid-cols-3',
        ]
    },
    // 🔥 @nuxt/image 配置 - 自动 WebP/AVIF 转换
    image: {
        format: ['webp', 'avif'],
        quality: 80,
        screens: {
            xs: 320,
            sm: 640,
            md: 768,
            lg: 1024,
            xl: 1280,
            xxl: 1536,
        },
        domains: [],  // 允许所有域名
        alias: {
            avatar: '/avatar',
        },
    },
    vue: {
        compilerOptions: {
            isCustomElement: (tag:string) => ['meting-js'].includes(tag),
        },
    },
    app: {
        head: {
            meta: [
                { name: "viewport", content: "width=device-width, initial-scale=1, user-scalable=no" },
                { charset: "utf-8" },
            ],
            link: [
                {href: `/css/APlayer.min.css`, rel: 'stylesheet'},
            ],
            script: [
                {src: `/js/APlayer.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/Meting.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/main.js`, type: 'text/javascript', async: true, defer: true},
            ]
        }
    },
    vite: {
        server: {
            proxy: {
                "/api": {
                    target: "http://localhost:37892",
                },
                "/upload": {
                    target: "http://localhost:37892",
                },
                "/rss": {
                    target: "http://localhost:37892",
                },
                "/swagger": {
                    target: "http://localhost:37892",
                },
            },
        },
        build: {
            rollupOptions: {
                output: {
                    hashCharacters: 'base36'
                }
            }
        }
    }
})