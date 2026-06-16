FROM cypress/included:15.16.0

WORKDIR /root



ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --frozen-lockfile --ignore-scripts
