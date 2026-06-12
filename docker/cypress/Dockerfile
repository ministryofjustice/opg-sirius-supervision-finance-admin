FROM cypress/included:15.17.0

WORKDIR /root

COPY package.json package-lock.json ./
COPY cypress/package.json cypress/package.json
RUN npm ci --workspace cypress

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
