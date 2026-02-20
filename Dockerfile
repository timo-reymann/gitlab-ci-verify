FROM busybox AS bin
COPY ./dist /binaries
RUN if [[ "$(arch)" == "x86_64" ]]; then \
        architecture="amd64"; \
    else \
        architecture="arm64"; \
    fi; \
    cp /binaries/gitlab-ci-verify_linux-${architecture} /bin/gitlab-ci-verify && \
    chmod +x /bin/gitlab-ci-verify && \
    chown 1000:1000 /bin/gitlab-ci-verify

FROM scratch as license
COPY LICENSE LICENSE
COPY NOTICE NOTICE

FROM chainguard/wolfi-base
LABEL org.opencontainers.image.title="gitlab-ci-verify"
LABEL org.opencontainers.image.description="Validate and lint your gitlab ci files using ShellCheck, the Gitlab API and curated checks"
LABEL org.opencontainers.image.ref.name="main"
LABEL org.opencontainers.image.licenses='GPL-3.0'
LABEL org.opencontainers.image.vendor="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.authors="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.url="https://github.com/timo-reymann/gitlab-ci-verify"
LABEL org.opencontainers.image.documentation="https://github.com/timo-reymann/gitlab-ci-verify"
LABEL org.opencontainers.image.source="https://github.com/timo-reymann/gitlab-ci-verify.git"

COPY --from=license / /

RUN apk add --no-cache bash \
    && adduser -D -u 1000 gitlab-ci-verify

ARG BUILD_TIME
ARG BUILD_VERSION
ARG BUILD_COMMIT_REF
LABEL org.opencontainers.image.created=$BUILD_TIME
LABEL org.opencontainers.image.version=$BUILD_VERSION
LABEL org.opencontainers.image.revision=$BUILD_COMMIT_REF

COPY --from=bin /bin/gitlab-ci-verify /bin/gitlab-ci-verify
WORKDIR /workspace
ENTRYPOINT [ "/bin/gitlab-ci-verify" ]