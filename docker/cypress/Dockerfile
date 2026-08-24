FROM cypress/included:15.20.1@sha256:327d62ee6d1c9bd7b4b490bb2772c1b70d37d56647953a4ec3105a939b7c61ff

WORKDIR /root



ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --frozen-lockfile --ignore-scripts
