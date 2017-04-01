import React, { useState, useEffect, useRef } from 'react';
import Layout from '@theme/Layout';
import Translate, { translate } from '@docusaurus/Translate';
import Link from '@docusaurus/Link';
import {
  Zap,
  Shield,
  Database,
  Network,
  ArrowRight,
  Code,
  Lock,
  Rocket,
  ShieldCheck,
  Tags,
  Globe
} from 'lucide-react';

// Intersection Observer Hook for scroll animations
const useIntersectionObserver = (options = {}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting) {
        setIsVisible(true);
        observer.disconnect();
      }
    }, { threshold: 0.1, ...options });

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => observer.disconnect();
  }, []);

  return [ref, isVisible] as const;
};

// Animated Section Wrapper
const AnimatedSection = ({ children, className = '', delay = 0 }: { children: React.ReactNode, className?: string, delay?: number }) => {
  const [ref, isVisible] = useIntersectionObserver();

  return (
    <div
      ref={ref}
      className={`transition-all duration-700 ease-out ${className} ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {children}
    </div>
  );
};

// Helper Components
const FeatureCard = ({ icon, title, desc, index }: { icon: React.ReactNode, title: string, desc: string, index: number }) => {
  const [ref, isVisible] = useIntersectionObserver();

  return (
    <div
      ref={ref}
      className={`p-6 rounded-xl bg-[#161b22] border border-slate-800 hover:border-green-500/50 hover:shadow-[0_0_20px_rgba(34,197,94,0.1)] transition-all duration-500 group h-full hover-lift ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
      style={{ transitionDelay: `${index * 100}ms` }}
    >
      <div className="w-12 h-12 bg-slate-800 rounded-lg flex items-center justify-center text-green-500 mb-4 group-hover:scale-110 group-hover:bg-green-500/20 transition-all duration-300">
        {icon}
      </div>
      <h3 className="text-xl font-bold text-white mb-3">{title}</h3>
      <p className="text-slate-400 leading-relaxed text-sm">{desc}</p>
    </div>
  );
};

const UseCaseCard = ({ icon, title, desc, index }: { icon: React.ReactNode, title: string, desc: string, index: number }) => {
  const [ref, isVisible] = useIntersectionObserver();

  return (
    <div
      ref={ref}
      className={`flex gap-4 transition-all duration-500 ${isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'}`}
      style={{ transitionDelay: `${index * 150}ms` }}
    >
      <div className="flex-shrink-0 mt-1">
        <div className="w-10 h-10 rounded-full bg-green-500/20 flex items-center justify-center text-green-500 group-hover:bg-green-500/30 transition-colors">
          {icon}
        </div>
      </div>
      <div>
        <h3 className="text-lg font-bold text-white mb-2">{title}</h3>
        <p className="text-slate-400 text-sm leading-relaxed">{desc}</p>
      </div>
    </div>
  );
};

// Tech Partner Logo Item
const TechLogo = ({ logo, name, index }: { logo: React.ReactNode, name: string, index: number }) => {
  const [ref, isVisible] = useIntersectionObserver();

  return (
    <div
      ref={ref}
      className={`flex items-center gap-2 sm:gap-3 text-lg sm:text-xl font-bold text-white logo-hover cursor-pointer transition-all duration-500 ${isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-90'}`}
      style={{ transitionDelay: `${index * 100}ms` }}
    >
      {logo}
      <span>{name}</span>
    </div>
  );
};

export default function Home(): React.ReactElement {
  const [terminalLines, setTerminalLines] = useState<number[]>([]);

  // Terminal typing animation
  useEffect(() => {
    const delays = [0, 300, 600, 900, 1200, 1500, 1800];
    delays.forEach((delay, index) => {
      setTimeout(() => {
        setTerminalLines(prev => [...prev, index]);
      }, delay);
    });
  }, []);

  return (
    <Layout
      title={translate({ id: 'homepage.title', message: 'The Open Source Hub for AI Models' })}
      description={translate({ id: 'homepage.description', message: 'World-class, self-hosted model repository designed for vLLM and SGLang. Secure your AI assets with RBAC, accelerate distribution with proxy caching.' })}
    >
      <main className="bg-[#0d1117] text-slate-300 font-sans selection:bg-green-500/30 selection:text-green-200 w-full overflow-x-hidden">

        {/* Hero Section */}
        <section className="relative pt-20 pb-16 md:pt-28 md:pb-24 lg:pt-36 lg:pb-32 overflow-hidden">
          {/* Animated Background Element */}
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[600px] md:w-[800px] lg:w-[1000px] h-[300px] md:h-[400px] lg:h-[500px] bg-green-500/10 blur-[100px] lg:blur-[120px] rounded-full pointer-events-none animate-pulse"></div>

          {/* Floating particles */}
          <div className="absolute inset-0 overflow-hidden pointer-events-none">
            <div className="absolute top-1/4 left-1/4 w-2 h-2 bg-green-500/30 rounded-full animate-float" style={{ animationDelay: '0s' }}></div>
            <div className="absolute top-1/3 right-1/4 w-3 h-3 bg-emerald-500/20 rounded-full animate-float" style={{ animationDelay: '1s' }}></div>
            <div className="absolute bottom-1/3 left-1/3 w-2 h-2 bg-green-400/25 rounded-full animate-float" style={{ animationDelay: '2s' }}></div>
          </div>

          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-12 items-center">
              {/* Left: Text Content */}
              <div className="text-center lg:text-left">
                <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-green-900/20 border border-green-800/50 text-green-400 text-xs font-semibold uppercase tracking-wider mb-6 animate-fade-in-up">
                  <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                  <Translate id="homepage.hero.comingSoon">Coming Soon</Translate>
                </div>
                <h1 className="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-extrabold text-white leading-tight mb-6 animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
                  <Translate
                    id="homepage.hero.title"
                    values={{
                      br: <br className="hidden sm:block" />,
                      hub: <span className="text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-emerald-600 animate-glow inline-block">Hub</span>
                    }}>
                    {'The Open Source {br}{hub} for AI Models'}
                  </Translate>
                </h1>
                <p className="text-base sm:text-lg text-slate-400 mb-8 max-w-xl mx-auto lg:mx-0 leading-relaxed animate-fade-in-up" style={{ animationDelay: '0.2s' }}>
                  <Translate 
                    id="homepage.hero.subtitle"
                    values={{
                      hf: <span className="text-white">Hugging Face</span>,
                      vllm: <span className="text-white">vLLM</span>,
                      sglang: <span className="text-white">SGLang</span>
                    }}>
                    {'MatrixHub is an open-source, self-hosted AI model registry engineered for large-scale enterprise inference. It serves as a drop-in private replacement for {hf}, purpose-built to accelerate {vllm} and {sglang} workloads.'}
                  </Translate>
                </p>

                <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start animate-fade-in-up" style={{ animationDelay: '0.3s' }}>
                  <Link
                    to="/docs/getting-started"
                    className="homepage-quick-start px-6 sm:px-8 py-3 sm:py-4 bg-green-600 hover:bg-green-500 rounded-lg font-bold text-base sm:text-lg transition-all shadow-[0_0_20px_rgba(22,163,74,0.4)] hover:shadow-[0_0_30px_rgba(22,163,74,0.6)] flex items-center justify-center gap-2 hover:no-underline hover:scale-105 transform"
                  >
                    <Translate id="homepage.hero.quickStart">Quick Start</Translate> <ArrowRight size={20} />
                  </Link>
                  <a
                    href="https://github.com/matrixhub-ai/matrixhub"
                    className="px-6 sm:px-8 py-3 sm:py-4 bg-[#161b22] hover:bg-[#21262d] text-white border border-slate-700 rounded-lg font-bold text-base sm:text-lg transition-all flex items-center justify-center gap-2 hover:text-white hover:no-underline hover:border-slate-600"
                  >
                    <Code size={20} /> <Translate id="homepage.hero.github">View on GitHub</Translate>
                  </a>
                </div>
              </div>

              {/* Right: Terminal / Visual */}
              <div className="relative group mt-8 lg:mt-0 animate-fade-in-right" style={{ animationDelay: '0.4s' }}>
                <div className="absolute -inset-1 bg-gradient-to-r from-green-600 to-emerald-600 rounded-xl blur opacity-25 group-hover:opacity-40 transition duration-1000"></div>
                <div className="relative rounded-xl bg-[#0d1117] border border-slate-700 shadow-2xl overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-3 bg-[#161b22] border-b border-slate-800">
                    <div className="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors cursor-pointer"></div>
                    <div className="w-3 h-3 rounded-full bg-yellow-500 hover:bg-yellow-400 transition-colors cursor-pointer"></div>
                    <div className="w-3 h-3 rounded-full bg-green-500 hover:bg-green-400 transition-colors cursor-pointer"></div>
                    <span className="ml-2 text-xs text-slate-500 font-mono">matrixhub-ai/matrixhub</span>
                  </div>
                  <div className="p-4 sm:p-6 font-mono text-xs sm:text-sm overflow-hidden min-h-[240px]">
                    {terminalLines.includes(0) && (
                      <div className="flex items-start text-slate-400 mb-2 animate-fade-in-left">
                        <span className="text-green-500 mr-2 flex-shrink-0">➜</span>
                        <span className="break-all">~ export HF_ENDPOINT=https://hub.matrix.internal</span>
                      </div>
                    )}
                    {terminalLines.includes(1) && (
                      <div className="flex items-start text-slate-400 mb-2 animate-fade-in-left">
                        <span className="text-green-500 mr-2 flex-shrink-0">➜</span>
                        <span className="break-all">~ vllm serve "deepseek-ai/deepseek-coder-33b"</span>
                      </div>
                    )}
                    {terminalLines.includes(2) && (
                      <div className="mt-4 text-slate-300 space-y-1">
                        {terminalLines.includes(3) && (
                          <div className="animate-fade-in-left"><span className="text-blue-400">INFO</span>: Connected to MatrixHub Proxy</div>
                        )}
                        {terminalLines.includes(4) && (
                          <div className="animate-fade-in-left"><span className="text-blue-400">INFO</span>: Cache HIT for layer 1 (sha256:a1b2...)</div>
                        )}
                        {terminalLines.includes(5) && (
                          <div className="animate-fade-in-left"><span className="text-blue-400">INFO</span>: Cache HIT for layer 2 (sha256:c3d4...)</div>
                        )}
                        {terminalLines.includes(6) && (
                          <div className="animate-fade-in-left"><span className="text-green-400">SUCCESS</span>: Model loaded in 2.4s (via LAN)</div>
                        )}
                      </div>
                    )}
                    <div className="mt-4 flex items-center gap-2">
                       <span className="text-green-500 animate-pulse">➜</span>
                       <span className="w-2 h-5 bg-slate-500 terminal-cursor"></span>
                    </div>
                  </div>
                </div>

                {/* Float Card - Hidden on mobile */}
                <div className="absolute -bottom-6 -right-6 bg-[#1c2128] p-4 rounded-lg border border-slate-700 shadow-xl hidden lg:block animate-float">
                  <div className="flex items-center gap-3">
                    <div className="p-2 bg-green-900/30 rounded-full text-green-400 animate-glow flex items-center justify-center">
                       <Zap size={20} />
                    </div>
                    <div>
                      <div className="text-xs text-slate-400"><Translate id="homepage.hero.downloadSpeed">Download Speed</Translate></div>
                      <div className="text-lg font-bold text-white">25.8 GB/s <span className="text-xs font-normal text-green-500"><Translate id="homepage.hero.intranet">Intranet</Translate></span></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Analogy Section */}
        <AnimatedSection>
          <section className="py-12 sm:py-16 border-y border-slate-800 bg-[#161b22]/50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
                <h2 className="text-xl sm:text-2xl font-medium text-slate-300 leading-relaxed">
                  <Translate
                    id="homepage.analogy.title"
                    values={{
                      matrixhub: <span className="text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-emerald-600 font-extrabold text-2xl sm:text-3xl">MatrixHub</span>,
                      hf: <span className="text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 to-orange-500 font-extrabold text-2xl sm:text-3xl">Hugging Face</span>,
                      harbor: <span className="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-500 font-extrabold text-2xl sm:text-3xl">Harbor</span>,
                      docker: <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-500 font-extrabold text-2xl sm:text-3xl">Docker Hub</span>
                    }}>
                    {'{matrixhub} is to {hf} what {harbor} is to {docker}.'}
                  </Translate>
                </h2>
               <p className="mt-4 text-slate-400 text-sm sm:text-base">
                 <Translate id="homepage.analogy.subtitle">Stop relying on public internet for mission-critical AI. Control your assets, accelerate your pipelines.</Translate>
               </p>
            </div>
          </section>
        </AnimatedSection>

        {/* Key Features */}
        <section id="features" className="py-16 sm:py-20 lg:py-24 bg-[#0d1117] relative">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <AnimatedSection className="text-center mb-12 sm:mb-16">
              <h2 className="text-green-500 font-bold tracking-wider uppercase text-xs sm:text-sm mb-2"><Translate id="homepage.features.label">Core Features</Translate></h2>
              <h3 className="text-2xl sm:text-3xl md:text-4xl font-bold text-white mb-4"><Translate id="homepage.features.title">Infrastructure designed for Scale</Translate></h3>
              <p className="text-slate-400 max-w-2xl mx-auto text-sm sm:text-base">
                <Translate id="homepage.features.subtitle">Built for the specific needs of SREs and Algorithm Engineers managing massive model weights.</Translate>
              </p>
            </AnimatedSection>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
              <FeatureCard
                icon={<Network size={24} />}
                title={translate({ id: 'homepage.features.1.title', message: 'Transparent HF Proxy' })}
                desc={translate({ id: 'homepage.features.1.desc', message: 'Drop-in replacement for Hugging Face. Point your HF_ENDPOINT to MatrixHub and keep all training/inference code unchanged.' })}
                index={0}
              />
              <FeatureCard
                icon={<Zap size={24} />}
                title={translate({ id: 'homepage.features.2.title', message: 'On-Demand Caching' })}
                desc={translate({ id: 'homepage.features.2.desc', message: 'Pull once, cache forever. Automatically localizes public models to slash redundant traffic and accelerate cluster-wide distribution.' })}
                index={1}
              />
              <FeatureCard
                icon={<Shield size={24} />}
                title={translate({ id: 'homepage.features.3.title', message: 'RBAC & Audit Logs' })}
                desc={translate({ id: 'homepage.features.3.desc', message: 'Fine-grained permissions, project-based isolation, and comprehensive audit trails for every upload and download.' })}
                index={2}
              />
              <FeatureCard
                icon={<Database size={24} />}
                title={translate({ id: 'homepage.features.4.title', message: 'Storage Agnostic' })}
                desc={translate({ id: 'homepage.features.4.desc', message: 'Compatible with local filesystems, NFS, and S3-compatible backends (MinIO, AWS). Scale to unlimited model capacity.' })}
                index={3}
              />
            </div>
          </div>
        </section>

        {/* Use Cases */}
