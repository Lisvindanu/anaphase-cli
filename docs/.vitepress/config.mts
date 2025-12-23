import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Anaphase',
  description: 'AI-Powered Golang Microservice Generator',

  locales: {
    root: {
      label: 'English',
      lang: 'en-US',
      link: '/'
    },
    id: {
      label: 'Bahasa Indonesia',
      lang: 'id-ID',
      link: '/id/',
      title: 'Anaphase',
      description: 'Generator Microservice Golang dengan AI'
    }
  },

  // SEO
  head: [
    // Favicon
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],

    // SEO Meta Tags
    ['meta', { name: 'keywords', content: 'golang, go, microservices, ddd, clean architecture, domain-driven design, code generator, ai, gemini, scaffolding, cli, tool' }],
    ['meta', { name: 'author', content: 'Anaphase' }],
    ['meta', { name: 'robots', content: 'index, follow' }],

    // Open Graph / Facebook
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:url', content: 'https://anaphygon.my.id/' }],
    ['meta', { property: 'og:title', content: 'Anaphase - AI-Powered Golang Microservice Generator' }],
    ['meta', { property: 'og:description', content: 'Generate production-ready Golang microservices with AI. From idea to deployment in minutes using Domain-Driven Design and Clean Architecture.' }],
    ['meta', { property: 'og:image', content: 'https://anaphygon.my.id/hero-image.svg' }],
    ['meta', { property: 'og:site_name', content: 'Anaphase' }],
    ['meta', { property: 'og:locale', content: 'en_US' }],

    // Twitter Card
    ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
    ['meta', { name: 'twitter:url', content: 'https://anaphygon.my.id/' }],
    ['meta', { name: 'twitter:title', content: 'Anaphase - AI-Powered Golang Microservice Generator' }],
    ['meta', { name: 'twitter:description', content: 'Generate production-ready Golang microservices with AI. From idea to deployment in minutes.' }],
    ['meta', { name: 'twitter:image', content: 'https://anaphygon.my.id/hero-image.svg' }],

    // Additional Meta
    ['meta', { name: 'theme-color', content: '#667eea' }],
    ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
    ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' }],

    // Canonical URL
    ['link', { rel: 'canonical', href: 'https://anaphygon.my.id/' }]
  ],

  themeConfig: {
    logo: '/logo.svg',

    // Search feature (Ctrl+K)
    search: {
      provider: 'local',
      options: {
        locales: {
          root: {
            translations: {
              button: {
                buttonText: 'Search',
                buttonAriaLabel: 'Search documentation'
              },
              modal: {
                noResultsText: 'No results for',
                resetButtonTitle: 'Clear search',
                footer: {
                  selectText: 'to select',
                  navigateText: 'to navigate',
                  closeText: 'to close'
                }
              }
            }
          },
          id: {
            translations: {
              button: {
                buttonText: 'Cari',
                buttonAriaLabel: 'Cari dokumentasi'
              },
              modal: {
                noResultsText: 'Tidak ada hasil untuk',
                resetButtonTitle: 'Hapus pencarian',
                footer: {
                  selectText: 'untuk memilih',
                  navigateText: 'untuk navigasi',
                  closeText: 'untuk menutup'
                }
              }
            }
          }
        }
      }
    },

    nav: [
      { text: 'Guide', link: '/guide/introduction' },
      { text: 'Quick Start', link: '/guide/quick-start' },
      { text: 'Reference', link: '/reference/commands' },
      { text: 'Examples', link: '/examples/basic' }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'What is Anaphase?', link: '/guide/introduction' },
          { text: 'Quick Start', link: '/guide/quick-start' },
          { text: 'Installation', link: '/guide/installation' },
          { text: '🆘 Troubleshooting', link: '/guide/troubleshooting' }
        ]
      },
      {
        text: 'Core Concepts',
        items: [
          { text: 'Architecture', link: '/guide/architecture' },
          { text: 'AI-Powered Generation', link: '/guide/ai-generation' },
          { text: 'Domain-Driven Design', link: '/guide/ddd' }
        ]
      },
      {
        text: 'Commands',
        items: [
          { text: 'anaphase init', link: '/reference/init' },
          { text: 'anaphase config', link: '/reference/config' },
          { text: 'anaphase gen domain', link: '/reference/gen-domain' },
          { text: 'anaphase gen handler', link: '/reference/gen-handler' },
          { text: 'anaphase gen repository', link: '/reference/gen-repository' },
          { text: 'anaphase gen middleware', link: '/reference/gen-middleware' },
          { text: 'anaphase gen migration', link: '/reference/gen-migration' },
          { text: 'anaphase quality', link: '/reference/quality' },
          { text: 'anaphase wire', link: '/reference/wire' }
        ]
      },
      {
        text: 'Examples',
        items: [
          { text: 'Basic E-commerce', link: '/examples/basic' },
          { text: 'Multi-Domain Service', link: '/examples/multi-domain' },
          { text: 'Custom Handlers', link: '/examples/custom-handlers' }
        ]
      },
      {
        text: 'Configuration',
        items: [
          { text: 'AI Providers', link: '/config/ai-providers' },
          { text: 'Database Settings', link: '/config/database' },
          { text: 'Project Structure', link: '/config/project-structure' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/lisvindanu/anaphase-cli' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present Anaphase'
    }
  }
})
