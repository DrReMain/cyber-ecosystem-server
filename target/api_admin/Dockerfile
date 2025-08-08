FROM alpine:latest

ARG PROJECT=api_admin
ARG CONFIG_FILE=api_admin.yaml

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}

ENV TZ=Asia/Shanghai
RUN apk update --no-cache && apk add --no-cache tzdata

COPY ./target/${PROJECT}/${PROJECT} ./
COPY ./target/${PROJECT}/${CONFIG_FILE} ./etc/

EXPOSE 9000

ENTRYPOINT ["./api_admin", "-f", "./etc/api_admin.yaml"]
