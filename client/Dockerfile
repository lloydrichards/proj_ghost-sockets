FROM node:22-alpine as builder

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist


COPY package.json .
COPY vite.config.ts .

RUN npm install typescript

EXPOSE 8080
CMD ["npm", "run", "preview"]