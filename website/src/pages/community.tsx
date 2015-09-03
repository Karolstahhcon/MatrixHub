import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function Community(): React.ReactElement {
  const sections = [
    {
      emoji: '⭐',
      title: translate({ id: 'community.github.title', message: 'GitHub' }),
      description: translate({ id: 'community.github.description', message: 'Star and follow the MatrixHub repository for updates, issues, and releases.' }),
      linkText: translate({ id: 'community.github.linkText', message: 'Visit Repository →' }),
      href: 'https://github.com/matrixhub-ai/matrixhub',
    },
    {
      emoji: '🤝',
      title: translate({ id: 'community.contributing.title', message: 'Contributing' }),
      description: translate({ id: 'community.contributing.description', message: 'We welcome contributions of all kinds — code, documentation, bug reports, and feature requests.' }),
      linkText: translate({ id: 'community.contributing.linkText', message: 'Contribution Guide →' }),
      href: 'https://github.com/matrixhub-ai/matrixhub/blob/main/CONTRIBUTING.md',
    },
    {
      emoji: '💬',
      title: translate({ id: 'community.channels.title', message: 'Community Channels' }),
      description: translate({ id: 'community.channels.description', message: 'Join the #matrixhub channel on CNCF Slack to chat with maintainers and users.' }),
      linkText: translate({ id: 'community.channels.linkText', message: 'Join #matrixhub on CNCF Slack →' }),
      href: 'https://cloud-native.slack.com/archives/C0A8UKWR8HG',
    },
  ];

  return (
    <Layout
      title={translate({ id: 'community.title', message: 'Community' })}
      description={translate({
        id: 'community.description',
        message:
          'Join the MatrixHub community — contribute, discuss, and build together.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        {/* Hero */}
