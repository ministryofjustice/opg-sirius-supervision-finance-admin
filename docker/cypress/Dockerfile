FROM cypress/included:15.21.0@sha256:af15ca5f3b851e32360f37e1d446d6e6821fd11928c87a2fd7fad12e8495e63b

WORKDIR /root



ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --frozen-lockfile --ignore-scripts
