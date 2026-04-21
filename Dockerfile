FROM crpi-wcwep2opqtm0wq85.cn-shenzhen.personal.cr.aliyuncs.com/private_image_repositor/node:20.19.1-bookworm AS front
WORKDIR /app
RUN npm install -g pnpm@10.10.0
COPY front/package.json .
COPY front/pnpm-lock.yaml .
COPY front/pnpm-workspace.yaml .
RUN pnpm install
COPY front/. .
RUN pnpm run generate

FROM crpi-wcwep2opqtm0wq85.cn-shenzhen.personal.cr.aliyuncs.com/private_image_repositor/golang:1.23.3-alpine AS backend
ARG VERSION
ARG COMMIT_ID
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod download
COPY backend/. .
COPY --from=front /app/.output/public /app/public
RUN go build -tags prod -ldflags="-s -w -X main.version=${VERSION} -X main.commitId=${COMMIT_ID}" -o /app/moments

FROM crpi-wcwep2opqtm0wq85.cn-shenzhen.personal.cr.aliyuncs.com/private_image_repositor/alpine:latest
WORKDIR /app/data
ENV PORT=3000
ENV TZ=Asia/Shanghai
COPY --from=backend /app/moments /app/moments
RUN chmod +x /app/moments
EXPOSE 3000
CMD ["/app/moments"]
