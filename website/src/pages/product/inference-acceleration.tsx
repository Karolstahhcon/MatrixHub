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
