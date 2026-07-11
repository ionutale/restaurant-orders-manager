FROM node:26-alpine AS builder
ENV CI=true PNPM_CONFIG_MINIMUM_RELEASE_AGE=0
RUN npm install -g pnpm@latest
WORKDIR /build
COPY package.json pnpm-lock.yaml .npmrc ./
RUN pnpm install
COPY . .
RUN pnpm lint && pnpm build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build/build /usr/share/nginx/html
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
