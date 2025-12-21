import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Anaphase',
  description: 'AI-Powered Golang Microservice Generator',

  themeConfig: {
    logo: '/logo.svg',

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
          { text: 'Installation', link: '/guide/installation' }
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
          { text: 'anaphase gen domain', link: '/reference/gen-domain' },
          { text: 'anaphase gen handler', link: '/reference/gen-handler' },
          { text: 'anaphase gen repository', link: '/reference/gen-repository' },
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
  },

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }]
  ]
})
