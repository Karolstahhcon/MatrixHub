import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function RemoteSync(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'sync.title', message: 'Remote Sync' })}
      description={translate({
        id: 'sync.description',
        message: 'Synchronize models across data centers and regions.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="sync.heading">Remote Sync</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="sync.subtitle">
                Synchronize models across data centers and regions.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-16">
          <div className="text-center max-w-3xl mx-auto mb-12">
            <h2 className="text-3xl font-bold text-white mb-6">
              <Translate id="sync.global.title">Global Multi-Region Replication</Translate>
            </h2>
            <p className="text-slate-400 text-lg leading-relaxed">
              <Translate id="sync.global.desc">
                Automate asynchronous, resumable replication between data centers for high availability and low-latency local access. MatrixHub makes global collaboration operationally manageable.
              </Translate>
            </p>
          </div>
