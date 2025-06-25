
FROM fluent/fluentd:v1.16-1
USER root
RUN apk add --no-cache curl

RUN gem install fluent-plugin-opensearch --no-document
RUN gem install fluent-plugin-docker --no-document

USER fluent
