FROM node:26-alpine AS builder
ENV CI=true PNPM_CONFIG_MINIMUM_RELEASE_AGE=0
RUN npm install -g pnpm@latest
WORKDIR /build
COPY package.json pnpm-lock.yaml .npmrc ./
RUN pnpm install
COPY . .
RUN pnpm lint && pnpm build

FROM node:26-alpine
WORKDIR /app
COPY --from=builder /build/build ./build
COPY --from=builder /build/package.json ./
EXPOSE 3000
ENV NODE_ENV=production
CMD ["node", "build"]
