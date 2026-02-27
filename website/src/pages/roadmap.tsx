import React from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';

export default function Roadmap(): React.ReactElement {
  return (
    <Layout
      title={translate({ id: 'roadmap.title', message: 'Roadmap' })}
      description={translate({
