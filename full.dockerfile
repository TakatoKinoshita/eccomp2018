FROM takatokinoshita/eccomp2018:build AS build

FROM takatokinoshita/eccomp2018:base AS release
COPY --from=build --chown=nonroot:nonroot --chmod=711 /app/moon_mop/moon_mop /app/moon_sop/moon_sop ./
COPY ./schema ./schema
ENV EVAL_MODULE=moon_mop
CMD ["./wrapper"]