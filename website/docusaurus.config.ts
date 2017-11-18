import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'MatrixHub',
  tagline: 'The Open Source Hub for AI Models',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  markdown: {
    mermaid: true,
  },

  url: 'https://matrixhub.ai',
  baseUrl: '/',

  organizationName: 'matrixhub-ai',
  projectName: 'matrixhub',

  onBrokenLinks: 'throw',
  onBrokenAnchors: 'throw',
  onBrokenMarkdownLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'zh-CN'],
    localeConfigs: {
      en: {
        label: 'English',
      },
      'zh-CN': {
        label: '简体中文',
      },
    },
  },

  presets: [
    [
      'classic',
      {
        docs: {
          path: './docs',
          routeBasePath: 'docs',
          sidebarPath: './sidebars.ts',
        },
        blog: {
          routeBasePath: 'blog',
          showReadingTime: true,
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themes: ['@docusaurus/theme-mermaid'],

  themeConfig: {
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: true,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'MatrixHub',
      logo: {
        alt: 'MatrixHub Logo',
        src: 'img/matrixhub-icon-colorlight.png',
        style: { height: '32px', width: '32px', borderRadius: '6px' },
      },
      hideOnScroll: false,
      items: [
        {
          to: '/product',
          label: 'Product',
          position: 'left',
        },
        {
          to: '/docs/overview',
          label: 'Docs',
          position: 'left',
        },
        {
          to: '/blog',
          label: 'Blog',
          position: 'left',
        },
        {
          to: '/community',
          label: 'Community',
          position: 'left',
        },
        {
          type: 'localeDropdown',
          position: 'right',
        },
        {
          href: 'https://github.com/matrixhub-ai/matrixhub',
          label: 'GitHub',
          position: 'right',
          className: 'navbar-github-link',
        },
