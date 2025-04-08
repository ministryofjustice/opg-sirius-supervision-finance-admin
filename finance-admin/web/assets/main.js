import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx;
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}];

const resetFieldInputs = () => {
    htmx.findAll('[id$="-field-input"]').forEach(element => {
        htmx.addClass(element, "hide");
        const input = element.querySelector("input");
        if (input) {
            input.setAttribute("disabled", "true");
            input.removeAttribute("max");
        }
    });
}

const showFieldInput = (idName) => {
    document.querySelector(`#${idName}`).removeAttribute("disabled");
    htmx.removeClass(htmx.find(`#${idName}-field-input`), "hide")
}

const yesterday = () => {
    const date = new Date();
    date.setDate(date.getDate() - 1);
    return date.toISOString().split('T')[0];
}

const setMaxDate = (idName, date) => {
    document.getElementById(idName).setAttribute("max", date);
}

function disableUploadFormInputs() {
    document.querySelector('#pis-number').setAttribute("disabled", "true")
    document.querySelector('#upload-date').setAttribute("disabled", "true")
    document.querySelector('#file-upload').setAttribute("disabled", "true")
    document.querySelector('#email-field').setAttribute("disabled", "true")
}

const dateRangeRequired = ["AgedDebt", "ARPaidInvoice", "TotalReceipts", "BadDebtWriteOff", "InvoiceAdjustments", "UnappliedReceipts"];

// adding event listeners inside the onLoad function will ensure they are re-added to partial content when loaded back in
htmx.onLoad(content => {
    initAll();

    htmx.findAll(".moj-banner--success").forEach((element) => {
        element.addEventListener("click", () => htmx.addClass(htmx.find(".moj-banner--success"), "hide"));
    });

    if (document.getElementById('reports-type')) {
        htmx.find("#reports-type").addEventListener("change", () => {
            const reportTypeEl = document.getElementById('reports-type');
            const reportType = reportTypeEl.value;
            resetFieldInputs();
            document.querySelector("form").reset();
            reportTypeEl.value =  reportType;

            switch (reportType) {
                case "Journal":
                    showFieldInput("journal-types");
                    showFieldInput("date");
                    showFieldInput("email");
                    setMaxDate("date", yesterday());
                    break;
                case "Schedule":
                    showFieldInput("schedule-types");
                    showFieldInput("date");
                    showFieldInput("email");
                    setMaxDate("date", yesterday());
                    break;
                case "AccountsReceivable":
                    showFieldInput("account-types");
                    showFieldInput("email");
                    break;
                case "Debt":
                    showFieldInput("debt-types");
                    showFieldInput("email");
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
            resetFieldInputs();
            document.querySelector("form").reset();
            reportTypeEl.value =  reportType;
            subTypeEl.value =  subType;

            showFieldInput("account-types");
            showFieldInput("email");

            if (dateRangeRequired.includes(subType)) {
                showFieldInput("date-to");
                showFieldInput("date-from");
            }
        });
    }

    if (document.getElementById('reports-upload-type')) {
        htmx.findAll("#reports-upload-type").forEach((element) => {
            element.addEventListener("change", function() {
                const elements = document.querySelectorAll('[id$="-input"]');
                elements.forEach(element => {
                    htmx.addClass(element, 'hide');
                });
                disableUploadFormInputs();
                const form = document.querySelector('form');
                const reportUploadTypeSelect = document.getElementById('reports-upload-type');
                const reportUploadTypeSelectValue = reportUploadTypeSelect.value;

                form.reset();
                reportUploadTypeSelect.value =  reportUploadTypeSelectValue;

                switch (reportUploadTypeSelect.value) {
                    case "PAYMENTS_SUPERVISION_CHEQUE":
                        document.querySelector('#pis-number').removeAttribute("disabled")
                        htmx.removeClass(htmx.find("#pis-number-input"), "hide")
                    case "PAYMENTS_MOTO_CARD":
                    case "PAYMENTS_ONLINE_CARD":
                    case "PAYMENTS_OPG_BACS":
                    case "PAYMENTS_SUPERVISION_BACS":
                    case "SOP_UNALLOCATED":
                    case "DIRECT_DEBITS_COLLECTIONS":
                        document.querySelector('#upload-date').removeAttribute("disabled");
                        document.querySelector('#file-upload').removeAttribute("disabled");
                        document.querySelector('#email-field').removeAttribute("disabled");
                        htmx.removeClass(htmx.find("#upload-date-input"), "hide");
                        htmx.removeClass(htmx.find("#file-upload-input"), "hide");
                        htmx.removeClass(htmx.find("#email-field-input"), "hide");
                        break;
                    case "DEBT_CHASE":
                    case "DEPUTY_SCHEDULE":
                        document.querySelector('#file-upload').removeAttribute("disabled")
                        htmx.addClass(htmx.find("#upload-date-input"), "hide")
                        htmx.addClass(htmx.find("#pis-number-input"), "hide")
                        htmx.addClass(htmx.find("#email-field-input"), "hide")
                        htmx.removeClass(htmx.find("#file-upload-input"), "hide")
                        break;
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