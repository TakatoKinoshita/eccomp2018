FROM takatokinoshita/eccomp2018:build AS build

FROM gcr.io/distroless/cc-debian11 AS release
WORKDIR /home/nonroot
USER nonroot
COPY --from=build --chown=nonroot:nonroot /app/DB ./DB
COPY --from=build --chown=nonroot:nonroot /app/wrapper .