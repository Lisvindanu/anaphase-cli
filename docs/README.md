# Anaphase Documentation

This is the documentation website for Anaphase CLI, built with VitePress.

## Development

### Install Dependencies

```bash
cd docs
npm install
```

### Run Development Server

```bash
npm run docs:dev
```

Visit http://localhost:5173

### Build for Production

```bash
npm run docs:build
```

Output in `.vitepress/dist/`

### Preview Production Build

```bash
npm run docs:preview
```

## Deployment

### GitHub Pages

1. Build the site:
   ```bash
   npm run docs:build
   ```

2. Deploy to GitHub Pages:
   ```bash
   # From docs directory
   cd .vitepress/dist
   git init
   git add -A
   git commit -m 'Deploy documentation'
   git push -f git@github.com:lisvindanuu/anaphase-cli.git main:gh-pages
   ```

3. Enable GitHub Pages in repository settings (source: gh-pages branch)

### Vercel

1. Import repository to Vercel
2. Set build settings:
   - Build command: `npm run docs:build`
   - Output directory: `docs/.vitepress/dist`
   - Install command: `cd docs && npm install`

### Netlify

1. Import repository to Netlify
2. Set build settings:
   - Base directory: `docs`
   - Build command: `npm run docs:build`
   - Publish directory: `docs/.vitepress/dist`

### Custom Server

Build and serve:

```bash
npm run docs:build
npx http-server .vitepress/dist
```

Or use nginx:

```nginx
server {
    listen 80;
    server_name docs.anaphase.dev;

    root /var/www/anaphase-docs;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## Structure

```
docs/
├── .vitepress/
│   └── config.mts        # VitePress configuration
├── guide/                # User guides
│   ├── introduction.md
│   ├── quick-start.md
│   ├── installation.md
│   ├── architecture.md
│   ├── ai-generation.md
│   └── ddd.md
├── reference/            # Command reference
│   ├── commands.md
│   ├── init.md
│   ├── gen-domain.md
│   ├── gen-handler.md
│   ├── gen-repository.md
│   └── wire.md
├── examples/             # Examples
│   └── basic.md
├── config/               # Configuration docs
├── index.md              # Homepage
└── package.json
```

## Contributing

To add new documentation:

1. Create `.md` file in appropriate directory
2. Add to sidebar in `.vitepress/config.mts`
3. Use frontmatter for metadata
4. Test locally with `npm run docs:dev`

## License

MIT
