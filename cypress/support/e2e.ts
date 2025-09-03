import "cypress-axe";
import "cypress-failed-log";
import * as axe from "axe-core";

declare global {
    namespace Cypress {
        interface Chainable {
            checkAccessibility(): Chainable<JQuery<HTMLElement>>
            setUser(id: string): Chainable<JQuery<HTMLElement>>
            login(email: string): Chainable<void>
            loginAs(user: string): Chainable<void>
        }
    }
}

Cypress.Commands.add("checkAccessibility", () => {
    const terminalLog = (violations: axe.Result[]) => {
        cy.task(
            "log",
            `${violations.length} accessibility violation${violations.length === 1 ? "" : "s"
            } ${violations.length === 1 ? "was" : "were"} detected`,
        );
        const violationData = violations.map(
            ({
                 id, impact, description, nodes,
             }) => ({
                id,
                impact,
                description,
                nodes: nodes.length,
            }),
        );
        cy.task("table", violationData);
    };
    cy.injectAxe();
    cy.configureAxe({
        rules: [
            {id: "aria-allowed-attr", selector: "*:not(input[type='radio'][aria-expanded])"},
        ],
    })
    cy.checkA11y(null, null, terminalLog);
});

Cypress.Commands.add("setUser", (id: string) => {
    cy.setCookie("x-test-user-id", id);
});

Cypress.Commands.add("login", (email: string): void => {
    let url = window.location.href;

    if (!(url.indexOf("localhost") > -1) && !(url.indexOf("finance-admin") > -1)) {
        cy.visit("/uploads")
        cy.get('a[href*="oauth/login"]').click();

        // Strip the current url, and use the stripped url prepended with auth to stop cypress freaking out about the redirect
        let authUrl = 'auth.' + url.slice(url.indexOf("//") + 2);
        authUrl = authUrl.slice(0, authUrl.indexOf("/"));

        cy.origin(authUrl, ({ args: email }), (email) => {
            cy.get('input[name="email"]').clear();
            cy.get('input[name="email"]').type(email);
            cy.get('[type="submit"]').click();
        })
    }
});

Cypress.Commands.add("loginAs", (user: string): void => {
    const emails = {
        "Allocations User": "allocations@opgtest.com",
        "Case Manager": "case.manager@opgtest.com",
        "Finance Manager": "finance.manager@opgtest.com",
        "Finance Reporting User": "finance.reporting@opgtest.com",
        "Finance User Testing": "finance.user.testing@opgtest.com",
        "LPA Manager": "2manager@opgtest.com",
        "Lay User": "Lay1-14@opgtest.com",
        "System Admin": "system.admin@opgtest.com",
        "Public API": "publicapi@opgtest.com",
    };

    const email = emails[user];

    if (email == null) {
        throw new Error("Could not find test login details for user " + user);
    }

    cy.login(email);
});
