FROM cypress/included:15.17.0@sha256:071a50a3acbef42b7070a19be7dd857b996f4e0ef154469059b99ab33c59bf53

WORKDIR /root



ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --frozen-lockfile --ignore-scripts
