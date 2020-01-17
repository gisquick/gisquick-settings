FROM node:12-alpine AS webapp
MAINTAINER Marcel Dancak "dancakm@gmail.com"

WORKDIR /webapp/
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build


FROM gisquick/webapp-container

COPY --from=webapp /webapp/dist/ /var/www
CMD /bin/true
