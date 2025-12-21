#!/bin/bash

# Build documentation
echo "Building documentation..."
npm run docs:build

# Navigate to build output
cd .vitepress/dist

# Initialize git if needed
if [ ! -d .git ]; then
  git init
  git remote add origin git@github.com:lisvindanuu/anaphase-cli.git
fi

# Deploy
echo "Deploying to GitHub Pages..."
git add -A
git commit -m "Deploy documentation - $(date '+%Y-%m-%d %H:%M:%S')"
git push -f origin main:gh-pages

echo "✅ Documentation deployed!"
echo "Visit: https://lisvindanuu.github.io/anaphase-cli/"
