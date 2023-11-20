FROM node:16-alpine3.15

RUN apk add curl
WORKDIR /app
COPY package.json .
RUN npm install

COPY . .

CMD [ "node", "index.js" ]
EXPOSE 3000