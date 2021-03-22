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
