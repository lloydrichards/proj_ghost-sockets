FROM node:22-alpine AS builder

ARG VITE_SERVER_HOST
ARG VITE_SERVER_PORT

ENV VITE_SERVER_HOST=$VITE_SERVER_HOST
ENV VITE_SERVER_PORT=$VITE_SERVER_PORT

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM node:22-alpine AS runner

WORKDIR /app

COPY --from=builder /app/dist ./dist


COPY package.json .
COPY vite.config.ts .

RUN npm install typescript

EXPOSE 8080
CMD ["npm", "run", "preview"]