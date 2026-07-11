FROM node:26-alpine AS builder
WORKDIR /build
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY . .
RUN corepack enable && pnpm lint && pnpm check && pnpm build

FROM node:26-alpine
WORKDIR /app
COPY --from=builder /build/build ./build
COPY --from=builder /build/package.json ./
EXPOSE 3000
ENV NODE_ENV=production
CMD ["node", "build"]
