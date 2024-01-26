FROM node:alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install -frozen-lockfile

COPY . .

RUN npm run build

ENV NODE_ENV=production

FROM node:alpine

ENV NODE_ENV=production

WORKDIR /app
USER 1000

COPY package*.json ./

COPY --from=builder /app/.next ./.next

COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules/ ./node_modules/


EXPOSE 3000

CMD ["npm", "start"]
