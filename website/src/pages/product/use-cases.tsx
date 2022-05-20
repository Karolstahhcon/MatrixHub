import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function UseCases(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'useCases.title', message: 'Use Cases' })}
      description={translate({
        id: 'useCases.description',
        message:
          'Real-world deployment scenarios for MatrixHub across enterprise AI teams.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="useCases.heading">Use Cases</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="useCases.subtitle">
                Real-world deployment scenarios for MatrixHub across enterprise AI teams.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-12">
          {/* Scenario 1 */}
          <div className="flex flex-col md:flex-row gap-8 items-start bg-slate-800/30 p-8 rounded-3xl border border-slate-700/50 hover:border-slate-600/50 transition-colors">
            <div className="md:w-1/3 shrink-0">
              <div className="text-sm font-bold text-green-500 mb-2 tracking-wider uppercase">
                <Translate id="useCases.uc1.label">Scenario 1</Translate>
              </div>
              <h2 className="text-2xl font-bold text-white mb-4">
                <Translate id="useCases.uc1.title">Intranet Inference Acceleration</Translate>
              </h2>
            </div>
            <div className="md:w-2/3 space-y-6">
              <div>
                <h3 className="text-lg font-semibold text-slate-200 mb-2">
                  <Translate id="useCases.uc1.scenarioLabel">The Scenario</Translate>
                </h3>
                <p className="text-slate-400 leading-relaxed">
                  <Translate id="useCases.uc1.scenario">An internal production environment runs a vLLM cluster with 100 GPU servers. A 70B model may exceed 130 GB. If every node pulls from the public Hugging Face endpoint independently, startup becomes slow and public bandwidth becomes a bottleneck.</Translate>
                </p>
              </div>
