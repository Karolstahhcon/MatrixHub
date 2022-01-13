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
