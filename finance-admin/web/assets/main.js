import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}]

function disableDownloadFormInputs() {
    document.querySelector('#report-journal-type').setAttribute("disabled", "true")
    document.querySelector('#report-schedule-type').setAttribute("disabled", "true")
    document.querySelector('#report-account-type').setAttribute("disabled", "true")
    document.querySelector('#report-debt-type').setAttribute("disabled", "true")
    document.querySelector('#date-field').setAttribute("disabled", "true")
    document.querySelector('#date-from-field').setAttribute("disabled", "true")
    document.querySelector('#date-to-field').setAttribute("disabled", "true")
    document.querySelector('#email-field').setAttribute("disabled", "true")
}

function disableUploadFormInputs() {
    document.querySelector('#upload-date').setAttribute("disabled", "true")
    document.querySelector('#file-upload').setAttribute("disabled", "true")
    document.querySelector('#email-field').setAttribute("disabled", "true")
}

// adding event listeners inside the onLoad function will ensure they are re-added to partial content when loaded back in
htmx.onLoad(content => {
    initAll();

    htmx.findAll(".moj-banner--success").forEach((element) => {
        element.addEventListener("click", () => htmx.addClass(htmx.find(".moj-banner--success"), "hide"));
    });

    if (document.getElementById('reports-type')) {
        htmx.findAll("#reports-type").forEach((element) => {
            element.addEventListener("change", function() {
                const elements = document.querySelectorAll('[id$="-field-input"]');
                elements.forEach(element => {
                    htmx.addClass(element, 'hide');
                });
                disableDownloadFormInputs();
                const form = document.querySelector('form');
                const reportTypeSelect = document.getElementById('reports-type');
                const reportTypeSelectValue = reportTypeSelect.value

                form.reset();
                reportTypeSelect.value =  reportTypeSelectValue

                switch (reportTypeSelect.value) {
                    case "Journal":
                        document.querySelector('#report-journal-type').removeAttribute("disabled");
                        document.querySelector('#date-field').removeAttribute("disabled");
                        htmx.removeClass(htmx.find("#journal-types-field-input"), "hide")
                        htmx.removeClass(htmx.find("#date-field-input"), "hide")
                        break;
                    case "Schedule":
                        document.querySelector('#report-schedule-type').removeAttribute("disabled");
                        document.querySelector('#date-field').removeAttribute("disabled");
                        htmx.removeClass(htmx.find("#schedule-types-field-input"), "hide")
                        htmx.removeClass(htmx.find("#date-field-input"), "hide")
                        break;
                    case "AccountsReceivable":
                        document.querySelector('#report-account-type').removeAttribute("disabled")
                        htmx.removeClass(htmx.find("#account-types-field-input"), "hide")
                        break;
                    case "Debt":
                        document.querySelector('#report-debt-type').removeAttribute("disabled")
                        htmx.removeClass(htmx.find("#debt-types-field-input"), "hide")
                        break;
                    default:
                        break;
                }
            }, false)
        });

        document.getElementById('report-account-type').addEventListener('change', function () {
            const form = document.querySelector('form');
            const reportTypeSelect = document.getElementById('reports-type');
            const reportTypeSelectValue = reportTypeSelect.value
            const reportAccountTypeSelectValue = this.value
            disableDownloadFormInputs();
            document.querySelector('#report-account-type').removeAttribute("disabled");

            form.reset();
            reportTypeSelect.value =  reportTypeSelectValue
            this.value = reportAccountTypeSelectValue

            switch (this.value) {
                case "AgedDebt":
                    document.querySelector('#email-field').removeAttribute("disabled");
                    document.querySelector('#date-to-field').removeAttribute("disabled");
                    document.querySelector('#date-from-field').removeAttribute("disabled");
                    htmx.removeClass(htmx.find("#email-field-input"), "hide")
                    htmx.removeClass(htmx.find("#date-to-field-input"), "hide")
                    htmx.removeClass(htmx.find("#date-from-field-input"), "hide")
                    break;
                case "AgedDebtByCustomer":
                    document.querySelector('#email-field').removeAttribute("disabled");
                    htmx.removeClass(htmx.find("#email-field-input"), "hide")
                    break;
                case "UnappliedReceipts":
                case "CustomerAgeingBuckets":
                    document.querySelector('#date-field').removeAttribute("disabled");
                    htmx.addClass(htmx.find("#date-to-field-input"), "hide")
                    htmx.addClass(htmx.find("#email-field-input"), "hide")
                    htmx.addClass(htmx.find("#date-from-field-input"), "hide")
                    htmx.removeClass(htmx.find("#date-field-input"), "hide")
                    break;
                case "ARPaidInvoiceReport":
                case "PaidInvoiceTransactionLines":
                case "TotalReceiptsReport":
                case "BadDebtWriteOffReport":
                    document.querySelector('#date-to-field').removeAttribute("disabled");
                    document.querySelector('#email-field').removeAttribute("disabled");
                    document.querySelector('#date-from-field').removeAttribute("disabled");
                    htmx.addClass(htmx.find("#date-field-input"), "hide")
                    htmx.removeClass(htmx.find("#date-to-field-input"), "hide")
                    htmx.removeClass(htmx.find("#email-field-input"), "hide")
                    htmx.removeClass(htmx.find("#date-from-field-input"), "hide")
                    break;
                case "FeeAccrual":
                    htmx.addClass(htmx.find("#date-field-input"), "hide")
                    htmx.addClass(htmx.find("#date-to-field-input"), "hide")
                    htmx.addClass(htmx.find("#email-field-input"), "hide")
                    htmx.addClass(htmx.find("#date-from-field-input"), "hide")
                    break;
                default:
                    break;
            }
        })
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
                const reportUploadTypeSelectValue = reportUploadTypeSelect.value

                form.reset();
                reportUploadTypeSelect.value =  reportUploadTypeSelectValue

                switch (reportUploadTypeSelect.value) {
                    case "PAYMENTS_MOTO_CARD":
                    case "PAYMENTS_ONLINE_CARD":
                    case "PAYMENTS_OPG_BACS":
                    case "PAYMENTS_SUPERVISION_BACS":
                        document.querySelector('#upload-date').removeAttribute("disabled")
                        document.querySelector('#file-upload').removeAttribute("disabled")
                        document.querySelector('#email-field').removeAttribute("disabled")
                        htmx.removeClass(htmx.find("#upload-date-input"), "hide")
                        htmx.removeClass(htmx.find("#file-upload-input"), "hide")
                        htmx.removeClass(htmx.find("#email-field-input"), "hide")
                        break;
                    case "DEBT_CHASE":
                    case "DEPUTY_SCHEDULE":
                        document.querySelector('#file-upload').removeAttribute("disabled")
                        htmx.addClass(htmx.find("#upload-date-input"), "hide")
                        htmx.addClass(htmx.find("#email-field-input"), "hide")
                        htmx.removeClass(htmx.find("#file-upload-input"), "hide")
                        break;
                    default:
                        break;
                }
            }, false)
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