# syntax=docker/dockerfile:1

FROM node:17.3.0-alpine as build

WORKDIR /frontend
COPY ./frontend .
RUN yarn
RUN yarn build

EXPOSE 3030

CMD ["node", ""]