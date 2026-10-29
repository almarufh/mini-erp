FROM node:lts-alpine3.24 AS builder

WORKDIR /erp

COPY package*.json ./

RUN npm ci --legacy-peer-deps

COPY . .
RUN npm run build

FROM nginx:1.31.6-alpine3.24

WORKDIR /erp
COPY --from=builder /erp/dist/ .
COPY --from=builder /erp/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]