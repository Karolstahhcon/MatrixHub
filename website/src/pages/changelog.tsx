import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function Changelog(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'changelog.title', message: 'Changelog' })}
      description={translate({
        id: 'changelog.description',
        message: 'Release history and version notes for MatrixHub.',
      })}
    >
      <main className="bg-[#0d1117] text-slate-300 min-h-screen">
