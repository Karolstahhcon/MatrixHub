import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';
import Link from '@docusaurus/Link';

export default function ProductOverview(): React.ReactElement {
  const features = [
    {
      title: translate({ id: 'product.feature.why.title', message: 'Why MatrixHub' }),
      description: translate({ id: 'product.feature.why.desc', message: 'Why organizations choose MatrixHub as their private model registry.' }),
      href: '/product/why-matrixhub',
      emoji: '🎯',
    },
    {
      title: translate({ id: 'product.feature.usecases.title', message: 'Use Cases' }),
      description: translate({ id: 'product.feature.usecases.desc', message: 'Real-world deployment scenarios across enterprise AI teams.' }),
      href: '/product/use-cases',
      emoji: '💼',
    },
    {
      title: translate({ id: 'product.feature.arch.title', message: 'Architecture' }),
      description: translate({ id: 'product.feature.arch.desc', message: 'Proxy layer, caching engine, storage backends, and API surface.' }),
      href: '/product/architecture',
      emoji: '🏗️',
    },
    {
      title: translate({ id: 'product.feature.hf.title', message: 'Hugging Face Compatible' }),
      description: translate({ id: 'product.feature.hf.desc', message: 'Drop-in compatibility with the Hugging Face ecosystem.' }),
      href: '/product/huggingface-compatible',
      emoji: '🤗',
    },
    {
      title: translate({ id: 'product.feature.inference.title', message: 'Inference Acceleration' }),
      description: translate({ id: 'product.feature.inference.desc', message: 'Accelerate model distribution for vLLM, SGLang, and more.' }),
      href: '/product/inference-acceleration',
      emoji: '⚡',
    },
    {
      title: translate({ id: 'product.feature.airgap.title', message: 'Air-Gapped Delivery' }),
      description: translate({ id: 'product.feature.airgap.desc', message: 'Securely deliver models to isolated environments.' }),
      href: '/product/air-gapped-delivery',
      emoji: '🔒',
    },
    {
      title: translate({ id: 'product.feature.sync.title', message: 'Remote Sync' }),
      description: translate({ id: 'product.feature.sync.desc', message: 'Synchronize models across data centers and regions.' }),
      href: '/product/remote-sync',
      emoji: '🔄',
    },
    {
      title: translate({ id: 'product.feature.gov.title', message: 'Governance' }),
      description: translate({ id: 'product.feature.gov.desc', message: 'RBAC, audit trails, and compliance controls for enterprise AI.' }),
      href: '/product/governance',
      emoji: '🛡️',
    },
  ];

  return (
    <Layout
      title={translate({ id: 'product.title', message: 'Product' })}
      description={translate({
        id: 'product.description',
        message:
