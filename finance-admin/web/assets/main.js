import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx;
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}];

const formToggler = (suffix) => {
    return {
        resetAll: resetAll(suffix),
        show: show(suffix),
    }
}

const resetAll = (suffix) => () => {
    htmx.findAll(`[id$="-${suffix}"]`).forEach(element => {
        htmx.addClass(element, "hide");
        const input = element.querySelector("input");
        if (input) {
            input.setAttribute("disabled", "true");
            input.removeAttribute("max");
        }
    });
}

const show = (suffix) => (idName) => {
    document.querySelector(`#${idName}`).removeAttribute("disabled");
    htmx.removeClass(htmx.find(`#${idName}-${suffix}`), "hide")
}

const yesterday = () => {
    const date = new Date();
    date.setDate(date.getDate() - 1);
    return date.toISOString().split('T')[0];
}

const setMaxDate = (idName, date) => {
    document.getElementById(idName).setAttribute("max", date);
}

const showFromDate = ["ARPaidInvoice", "TotalReceipts", "BadDebtWriteOff", "InvoiceAdjustments", "UnappliedReceipts", "AllRefunds"];
const showToDate = ["AgedDebt", "ARPaidInvoice", "TotalReceipts", "BadDebtWriteOff", "InvoiceAdjustments", "UnappliedReceipts", "AllRefunds"];

// adding event listeners inside the onLoad function will ensure they are re-added to partial content when loaded back in
htmx.onLoad(content => {
    initAll();

    htmx.findAll(".moj-banner--success").forEach((element) => {
        element.addEventListener("click", () => htmx.addClass(htmx.find(".moj-banner--success"), "hide"));
    });

    if (document.getElementById('reports-type')) {
        const toggle = formToggler("field-input")
        htmx.find("#reports-type").addEventListener("change", () => {
            const reportTypeEl = document.getElementById('reports-type');
            const reportType = reportTypeEl.value;
            toggle.resetAll();
            document.querySelector("form").reset();
            reportTypeEl.value =  reportType;

            switch (reportType) {
                case "Journal":
                    toggle.show("journal-types");
                    toggle.show("date");
                    toggle.show("email");
                    setMaxDate("date", yesterday());
                    break;
                case "Schedule":
                    toggle.show("schedule-types");
                    toggle.show("date");
                    toggle.show("email");
                    setMaxDate("date", yesterday());
                    break;
                case "AccountsReceivable":
                    toggle.show("account-types");
                    toggle.show("email");
                    break;
                case "Debt":
                    toggle.show("debt-types");
                    toggle.show("email");
                    break;
                default:
                    break;
            }
        }, false);

        htmx.find("#account-types").addEventListener("change", () => {
            const reportTypeEl = document.getElementById('reports-type');
            const reportType = reportTypeEl.value;
            const subTypeEl = document.getElementById('account-types');
            const subType = subTypeEl.value;

            toggle.resetAll();
            document.querySelector("form").reset();
            reportTypeEl.value =  reportType;
            subTypeEl.value =  subType;

            toggle.show("account-types");
            toggle.show("email");

            if (showFromDate.includes(subType)) {
                toggle.show("date-from");
            }

            if (showToDate.includes(subType)) {
                toggle.show("date-to");
            }
        });

        htmx.find('#schedule-types').addEventListener("change", () => {
            const reportTypeEl = document.getElementById('reports-type');
            const reportType = reportTypeEl.value;
            const subTypeEl = document.getElementById('schedule-types');
            const subType = subTypeEl.value;

            toggle.resetAll();
            document.querySelector("form").reset();
            reportTypeEl.value =  reportType;
            subTypeEl.value =  subType;

            toggle.show("schedule-types");
            toggle.show("date");
            toggle.show("email");
            setMaxDate("date", yesterday());

            if (subType === 'ChequePayments') {
                toggle.show("pis-number")
            }
        });
    }

    if (document.getElementById('upload-type')) {
        htmx.findAll("#upload-type").forEach((element) => {
            const toggle = formToggler("input");
            element.addEventListener("change", function() {
                toggle.resetAll();
                const form = document.querySelector('form');
                const reportUploadTypeSelect = document.getElementById('upload-type');
                const reportUploadTypeSelectValue = reportUploadTypeSelect.value;
                form.reset();

                reportUploadTypeSelect.value =  reportUploadTypeSelectValue;

                switch (reportUploadTypeSelect.value) {
                    case "PAYMENTS_SUPERVISION_CHEQUE":
                        toggle.show("pis-number");
                        toggle.show("upload-date");
                        toggle.show("file-upload");
                        toggle.show("email-field");
                        break;
                    case "PAYMENTS_MOTO_CARD":
                    case "PAYMENTS_ONLINE_CARD":
                    case "PAYMENTS_OPG_BACS":
                    case "PAYMENTS_SUPERVISION_BACS":
                    case "DIRECT_DEBITS_COLLECTIONS":
                    case "FULFILLED_REFUNDS":
                    case "REVERSE_FULFILLED_REFUNDS":
                        toggle.show("upload-date");
                        toggle.show("file-upload");
                        toggle.show("email-field");
                        break;
                    case "DEBT_CHASE":
                    case "DEPUTY_SCHEDULE":
                    case "MISAPPLIED_PAYMENTS":
                    case "DUPLICATED_PAYMENTS":
                    case "BOUNCED_CHEQUE":
                    case "FAILED_DIRECT_DEBITS_COLLECTIONS":
                    case "SOP_UNALLOCATED":
                        toggle.show("file-upload");
                        toggle.show("email-field");
                        break
                    default:
                        break;
                }
            }, false);
        });
    }

    // validation errors are loaded in as a partial, with oob-swaps for the field error messages,
    // but classes need to be applied to each form group that appears in the summary
    const errorSummary = htmx.find("#error-summary");
    if (errorSummary) {
        const errors = [];
        errorSummary.querySelectorAll(".govuk-link").forEach((element) => {
            errors.push(element.getAttribute("href"));
        });
        htmx.findAll(".govuk-form-group").forEach((element) => {
            if (errors.includes(`#${element.id}`)) {
                element.classList.add("govuk-form-group--error");
                element.querySelector('.govuk-error-message')?.classList.remove('hide');
            } else {
                element.classList.remove("govuk-form-group--error");
                element.querySelector('.govuk-error-message')?.classList.add('hide');
            }
        })
    }
});