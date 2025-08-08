FROM alpine:latest

ARG PROJECT=rpc_admin_system
ARG CONFIG_FILE=rpc_admin_system.yaml

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}

ENV TZ=Asia/Shanghai
RUN apk update --no-cache && apk add --no-cache tzdata

COPY ./target/${PROJECT}/${PROJECT} ./
COPY ./target/${PROJECT}/${CONFIG_FILE} ./etc/
COPY ./target/data/menu.json ./etc/
COPY ./target/data/role.json ./etc/
COPY ./target/data/user.json ./etc/

EXPOSE 8001

ENTRYPOINT ["./rpc_admin_system", "-f", "./etc/rpc_admin_system.yaml"]
