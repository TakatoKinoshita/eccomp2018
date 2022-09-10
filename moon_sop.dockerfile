FROM golang:1.19-bullseye AS build
WORKDIR /app
RUN apt-get update && apt-get install -y unar build-essential aria2
RUN aria2c -x16 -s16 http://www.jpnsec.org/files/competition2018/data/DB.zip \
    && unar ./DB.zip
RUN wget http://www.jpnsec.org/files/competition2018/data/moon_sop.zip \
    && unar ./moon_sop.zip \
    && wget http://www.jpnsec.org/files/competition2018/data/moon_mop.tgz \
    && tar xvf ./moon_mop.tgz
RUN cd moon_sop && make
RUN cd moon_mop && make
COPY go.mod go.sum wrapper.go ./
RUN ["go", "build"]

FROM gcr.io/distroless/cc-debian11 AS release
WORKDIR /home/nonroot
USER nonroot
COPY --from=build --chown=nonroot:nonroot /app/DB ./DB
COPY --from=build --chown=nonroot:nonroot /app/wrapper .
COPY --from=build --chown=nonroot:nonroot --chmod=711 /app/moon_sop/moon_sop ./
COPY ./schema/moon_sop.json ./schema/
ENV EVAL_MODULE=moon_sop
CMD ["./wrapper"]