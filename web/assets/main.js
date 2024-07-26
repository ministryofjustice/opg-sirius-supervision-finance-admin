import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}]
