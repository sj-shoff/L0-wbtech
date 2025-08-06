FROM nginx:alpine

COPY index.html script.js /usr/share/nginx/html/

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d