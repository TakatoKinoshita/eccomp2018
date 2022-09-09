FROM golang:1.18.4-bullseye AS build
WORKDIR /app
RUN apt-get update && apt-get install -y unzip
RUN wget https://ladse.eng.isas.jaxa.jp/benchmark/Mazda_CdMOBP.zip \
    && unzip ./Mazda_CdMOBP.zip
COPY go.mod .
COPY go.sum .
COPY wrapper.go .
RUN ["go", "build"]

FROM gcr.io/distroless/cc-debian11 AS release
WORKDIR /home/nonroot
USER nonroot
COPY --from=build --chown=nonroot:nonroot --chmod=711 /app/Mazda_CdMOBP/Mazda_CdMOBP/bin/Linux/* .
COPY --from=build --chown=nonroot:nonroot /app/wrapper .
COPY ./schema ./schema
ENV EVAL_MODULE=mazda_mop
CMD ["./wrapper"]