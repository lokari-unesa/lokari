# syntax=docker/dockerfile:1

# Frontend Docker HANYA untuk DEVELOPMENT (vite dev server, hot reload).
# Deploy produksi frontend lewat Vercel (adapter-vercel) — lihat README.
FROM node:22-slim AS dev
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
EXPOSE 5180
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
