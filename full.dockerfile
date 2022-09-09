FROM golang:1.19-bullseye AS build
WORKDIR /app
RUN apt-get update && apt-get install -y unzip build-essential
RUN wget http://www.jpnsec.org/files/competition2018/data/moon_sop.zip \
    && unzip ./moon_sop.zip \
    && wget http://www.jpnsec.org/files/competition2018/data/moon_mop.tgz \
    && tar -xzvf ./moon_mop.tgz
RUN cd moon_sop && make
RUN cd moon_mop && make
RUN wget http://www.jpnsec.org/files/competition2018/data/DB.zip
RUN unzip ./DB.zip
COPY go.mod .
COPY go.sum .
COPY wrapper.go .
RUN ["go", "build"]

# FROM gcr.io/distroless/cc-debian11 AS release
# WORKDIR /home/nonroot
# USER nonroot
# COPY --from=build --chown=nonroot:nonroot --chmod=711 /app/Mazda_CdMOBP/Mazda_CdMOBP/bin/Linux/* .
# COPY --from=build --chown=nonroot:nonroot /app/wrapper .
# COPY ./schema ./schema
# ENV EVAL_MODULE=mazda_mop
# CMD ["./wrapper"]