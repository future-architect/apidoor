FROM node:16.7 as builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY tsconfig.json vue.config.js babel.config.js ./
COPY src ./src
RUN npm run build

FROM nginx:1.21-alpine as runner

COPY --from=builder /app/dist /public
COPY contrib/nginx.conf.template /etc/nginx/templates/default.conf.template
CMD ["nginx", "-g", "daemon off;"]
