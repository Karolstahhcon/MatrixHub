import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function WhyMatrixHub(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'whyMatrixhub.title', message: 'Why MatrixHub' })}
      description={translate({
        id: 'whyMatrixhub.description',
        message:
          'Why organizations choose MatrixHub as their private model registry for enterprise AI.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="whyMatrixhub.heading">Why MatrixHub</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="whyMatrixhub.subtitle">
                Why organizations choose MatrixHub as their private model registry for enterprise AI.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-20">
          <div>
            <h2 className="text-3xl font-bold text-white mb-8 border-b border-slate-800 pb-4">
              <Translate id="why.problem.title">The Challenge of Enterprise AI</Translate>
