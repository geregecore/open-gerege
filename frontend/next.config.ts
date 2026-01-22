import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  // Allow images from trusted domains
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '*.gerege.mn',
      },
      {
        protocol: 'https',
        hostname: 'sso.gerege.mn',
      },
      {
        protocol: 'https',
        hostname: 'core.gerege.mn',
      },
    ],
  },
  async rewrites() {
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';
    return [
      {
        source: '/api/v1/:path*',
        destination: `${backendUrl}/:path*`, // Proxy to backend, stripping /api/v1
      },
    ];
  },
};

export default nextConfig;
