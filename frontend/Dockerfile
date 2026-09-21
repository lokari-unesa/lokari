# syntax=docker/dockerfile:1

# ============ BASE — dependensi & kode sumber ============
FROM node:22-slim AS base
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .

# ============ DEV — vite dev server ============
FROM base AS dev
EXPOSE 5180
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# ============ BUILD — bundle produksi (adapter-node) ============
FROM base AS build
RUN npm run build

# ============ PROD — server SvelteKit adapter-node ============
FROM node:22-alpine AS prod
WORKDIR /app
ENV NODE_ENV=production HOST=0.0.0.0 PORT=5180
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./
COPY --from=build /app/node_modules ./node_modules
EXPOSE 5180
CMD ["node", "build"]
