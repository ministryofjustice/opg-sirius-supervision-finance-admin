FROM cypress/included:15.6.0

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe
RUN npm install aws-sdk

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://finance-admin:8888/finance-admin

COPY cypress.config.js .
COPY cypress cypress
