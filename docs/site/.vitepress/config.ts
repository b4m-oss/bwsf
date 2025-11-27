import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'bwenv',
  description: 'Manage .env files with Bitwarden CLI',
  
  // GitHub Pages base path
  base: '/bwenv/',
  
  // Clean URLs (no .html extension)
  cleanUrls: true,
  
  // Head meta tags
  head: [
    ['meta', { name: 'theme-color', content: '#175ddc' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:site_name', content: 'bwenv' }],
  ],
  
  // i18n configuration
  locales: {
    en: {
      label: 'English',
      lang: 'en',
      link: '/en/',
      themeConfig: {
        nav: [
          { text: 'Guide', link: '/en/guide/getting-started' },
          { text: 'Commands', link: '/en/guide/commands' },
          {
            text: 'Links',
            items: [
              { text: 'GitHub', link: 'https://github.com/b4m-oss/bwenv' },
              { text: 'Changelog', link: 'https://github.com/b4m-oss/bwenv/releases' },
            ]
          }
        ],
        sidebar: [
          {
            text: 'Introduction',
            items: [
              { text: 'What is bwenv?', link: '/en/guide/getting-started' },
              { text: 'Installation', link: '/en/guide/installation' },
            ]
          },
          {
            text: 'Usage',
            items: [
              { text: 'Commands', link: '/en/guide/commands' },
            ]
          }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © b4m-oss'
        },
        editLink: {
          pattern: 'https://github.com/b4m-oss/bwenv/edit/main/docs/site/:path',
          text: 'Edit this page on GitHub'
        }
      }
    },
    ja: {
      label: '日本語',
      lang: 'ja',
      link: '/ja/',
      themeConfig: {
        nav: [
          { text: 'ガイド', link: '/ja/guide/getting-started' },
          { text: 'コマンド', link: '/ja/guide/commands' },
          {
            text: 'リンク',
            items: [
              { text: 'GitHub', link: 'https://github.com/b4m-oss/bwenv' },
              { text: '変更履歴', link: 'https://github.com/b4m-oss/bwenv/releases' },
            ]
          }
        ],
        sidebar: [
          {
            text: 'はじめに',
            items: [
              { text: 'bwenvとは？', link: '/ja/guide/getting-started' },
              { text: 'インストール', link: '/ja/guide/installation' },
            ]
          },
          {
            text: '使い方',
            items: [
              { text: 'コマンド', link: '/ja/guide/commands' },
            ]
          }
        ],
        footer: {
          message: 'MITライセンスの下で公開されています。',
          copyright: 'Copyright © b4m-oss'
        },
        editLink: {
          pattern: 'https://github.com/b4m-oss/bwenv/edit/main/docs/site/:path',
          text: 'GitHubでこのページを編集'
        },
        docFooter: {
          prev: '前のページ',
          next: '次のページ'
        },
        outline: {
          label: '目次'
        },
        lastUpdated: {
          text: '最終更新'
        },
        returnToTopLabel: 'トップに戻る',
        sidebarMenuLabel: 'メニュー',
        darkModeSwitchLabel: 'ダークモード',
      }
    }
  },
  
  themeConfig: {
    // Logo and site title
    logo: '/logo.svg',
    siteTitle: 'bwenv',
    
    // Social links
    socialLinks: [
      { icon: 'github', link: 'https://github.com/b4m-oss/bwenv' }
    ],
    
    // Search
    search: {
      provider: 'local'
    }
  }
})
