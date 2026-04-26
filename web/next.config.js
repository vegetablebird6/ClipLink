const path = require('path');

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  distDir: 'dist',
  turbopack: {
    root: path.resolve(__dirname),
  },
};

module.exports = nextConfig; 
