import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function HuggingFaceCompatible(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'hf.title', message: 'Hugging Face Compatible' })}
      description={translate({
        id: 'hf.description',
        message: 'Drop-in compatibility with the Hugging Face ecosystem.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
            <h1 className="text-4xl sm:text-5xl font-extrabold text-white mb-6">
              <Translate id="hf.heading">Hugging Face Compatible</Translate>
            </h1>
            <p className="text-lg text-slate-400 max-w-2xl mx-auto">
              <Translate id="hf.subtitle">
                Drop-in compatibility with the Hugging Face ecosystem.
              </Translate>
            </p>
          </div>
        </section>

        <section className="py-16 max-w-5xl mx-auto px-6 space-y-16">
          <div className="text-center max-w-3xl mx-auto">
            <h2 className="text-3xl font-bold text-white mb-6">
              <Translate id="hf.transparent.title">Transparent Proxy Mode</Translate>
            </h2>
            <p className="text-slate-400 text-lg">
              <Translate id="hf.transparent.desc">
                Switch to private hosting with zero code changes. MatrixHub provides a drop-in endpoint replacement that your existing tools already understand.
              </Translate>
            </p>
          </div>

          <div className="bg-slate-900 rounded-2xl p-8 border border-slate-800 shadow-2xl relative overflow-hidden">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-green-500 to-emerald-400"></div>
            <div className="flex items-center justify-between mb-6 border-b border-slate-800 pb-4">
              <div className="text-slate-400 font-mono text-sm">Terminal</div>
              <div className="flex gap-2">
                <div className="w-3 h-3 rounded-full bg-slate-700"></div>
                <div className="w-3 h-3 rounded-full bg-slate-700"></div>
                <div className="w-3 h-3 rounded-full bg-slate-700"></div>
              </div>
            </div>
            <pre className="font-mono text-sm text-slate-300 overflow-x-auto">
              <code>
<span className="text-slate-500"># 1. Point your environment to MatrixHub</span><br/>
<span className="text-green-400">export</span> HF_ENDPOINT=<span className="text-yellow-300">"http://matrixhub.internal.company.com:3001"</span><br/>
<br/>
<span className="text-slate-500"># 2. Run your existing inference engine (e.g. vLLM)</span><br/>
<span className="text-green-400">python</span> -m vllm.entrypoints.openai.api_server \<br/>
&nbsp;&nbsp;--model deepseek-ai/DeepSeek-Coder-V2-Lite-Instruct<br/>
<br/>
<span className="text-slate-500"># MatrixHub automatically proxies the request, caches the model locally,</span><br/>
<span className="text-slate-500"># and serves it to all subsequent nodes at 10Gbps+ speeds.</span>
              </code>
