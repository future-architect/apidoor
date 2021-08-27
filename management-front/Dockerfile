FROM node:16.7

COPY package.json package-lock.json ./
RUN npm install
COPY tsconfig.json vue.config.js babel.config.js ./
COPY src ./src
EXPOSE 8080
CMD ["npm", "run", "serve"]
