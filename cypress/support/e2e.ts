import "cypress-axe";
import "cypress-failed-log";
import * as axe from "axe-core";

declare global {
    namespace Cypress {
        interface Chainable {
            checkAccessibility(): Chainable<JQuery<HTMLElement>>
            setUser(id: string): Chainable<JQuery<HTMLElement>>
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
