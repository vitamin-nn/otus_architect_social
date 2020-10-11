FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install --no-optional
COPY . .
ARG API_SERVER_URL=https://otus-architect.herokuapp.com/api/
RUN sed -i 's|$API_SERVER_URL|'"$API_SERVER_URL"'|g' .env.production
RUN echo $API_SERVER_URL > test.txt
RUN npm run build

FROM nginx:stable-alpine as production-stage
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY --from=build-stage /app/.env* /usr/share/nginx/html/
COPY --from=build-stage /app/test.txt /usr/share/nginx/html/
RUN echo $PORT > /usr/share/nginx/html/test-port.txt
CMD sed -i -e 's/$PORT/'"$PORT"'/g' /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'
