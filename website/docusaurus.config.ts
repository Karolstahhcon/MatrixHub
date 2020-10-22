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
