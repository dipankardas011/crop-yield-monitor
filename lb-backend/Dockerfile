FROM nginx:alpine

RUN rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/nginx.conf

COPY swaggerui /usr/share/nginx/html/swaggerui

EXPOSE 80
