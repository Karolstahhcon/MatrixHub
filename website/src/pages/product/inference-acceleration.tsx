import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function InferenceAcceleration(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'acceleration.title', message: 'Inference Acceleration' })}
      description={translate({
        id: 'acceleration.description',
        message: 'Accelerate model distribution for vLLM, SGLang, and more.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="acceleration.heading">Inference Acceleration</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="acceleration.subtitle">
                Accelerate model distribution for vLLM, SGLang, and more.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-16">
          <div className="text-center max-w-3xl mx-auto mb-12">
            <h2 className="text-3xl font-bold text-white mb-6">
              <Translate id="acceleration.zeroWait.title">Zero-Wait Distribution</Translate>
            </h2>
            <p className="text-slate-400 text-lg leading-relaxed">
              <Translate id="acceleration.zeroWait.desc">
                Eliminate bandwidth bottlenecks with a "Pull-once, serve-all" cache, enabling 10Gbps+ speeds across 100+ GPU nodes simultaneously. MatrixHub completely changes the startup performance of large inference clusters.
              </Translate>
            </p>
          </div>

