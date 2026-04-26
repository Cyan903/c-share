FROM node:lts

# Install
RUN npm install -g pnpm serve
WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN pnpm i

# Build
ARG VITE_APP
ARG VITE_API
ARG VITE_SOURCE

ENV VITE_APP=${VITE_APP}
ENV VITE_API=${VITE_API}
ENV VITE_SOURCE=${VITE_SOURCE}

COPY . .
RUN pnpm run build

# Host
CMD ["serve", "-s", "dist", "-l", "3000"]
