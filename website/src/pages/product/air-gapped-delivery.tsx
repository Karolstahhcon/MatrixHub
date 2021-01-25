import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function AirGappedDelivery(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'airgap.title', message: 'Air-Gapped Delivery' })}
      description={translate({
        id: 'airgap.description',
        message: 'Securely deliver models to isolated environments.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="airgap.heading">Air-Gapped Delivery</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="airgap.subtitle">
                Securely deliver models to isolated environments.
              </Translate>
            </p>
          </div>
        </section>

