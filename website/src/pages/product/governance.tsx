import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function Governance(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'governance.title', message: 'Governance' })}
      description={translate({
        id: 'governance.description',
        message: 'RBAC, audit trails, and compliance controls for enterprise AI.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="governance.heading">Governance</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="governance.subtitle">
                RBAC, audit trails, and compliance controls for enterprise AI.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-16">
          <div className="text-center max-w-3xl mx-auto mb-12">
            <h2 className="text-3xl font-bold text-white mb-6">
              <Translate id="governance.artifact.title">Enterprise Artifact Management</Translate>
            </h2>
            <p className="text-slate-400 text-lg leading-relaxed">
              <Translate id="governance.artifact.desc">
                Multiple internal teams fine-tune models. Versions drift between training, testing, and production. MatrixHub turns models into governed production artifacts instead of informal file shares.
              </Translate>
            </p>
          </div>

