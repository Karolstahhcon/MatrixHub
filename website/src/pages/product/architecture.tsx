import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function Architecture(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'architecture.title', message: 'Architecture' })}
      description={translate({
        id: 'architecture.description',
        message:
          'Technical architecture of MatrixHub — proxy layer, caching engine, storage backends, and API surface.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
        <section className="py-20 relative overflow-hidden">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-green-500/10 blur-[120px] rounded-full pointer-events-none"></div>
          <div className="max-w-4xl mx-auto px-6 relative z-10 text-center">
