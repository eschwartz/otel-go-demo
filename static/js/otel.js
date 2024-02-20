// Adapted from https://opentelemetry.io/docs/languages/js/getting-started/browser/
import {
    BatchSpanProcessor,
    ConsoleSpanExporter,
    SimpleSpanProcessor,
    WebTracerProvider
} from '@opentelemetry/sdk-trace-web';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import {OTLPTraceExporter} from "@opentelemetry/exporter-trace-otlp-http";
import {Resource} from "@opentelemetry/resources";
import {UserInteractionInstrumentation} from "@opentelemetry/instrumentation-user-interaction";
import {XMLHttpRequestInstrumentation} from "@opentelemetry/instrumentation-xml-http-request";
import {FetchInstrumentation} from "@opentelemetry/instrumentation-fetch";

const provider = new WebTracerProvider({
    resource: new Resource({
        // NOTE: this is used as the dataset in honeycomb
        "service.name": 'test',
    }),
});

// Configure exporter
const exporter = new OTLPTraceExporter({
    url: "https://api.honeycomb.io/v1/traces", // US instance
    headers: {
        // NOTE: As this is running client side, there is no way to hide this API key
        // A better alternative may be to run an otel collector on the same server
        "x-honeycomb-team": process.env.HONEYCOMB_API_KEY,
    },
})
provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

provider.register({
    // Changing default contextManager to use ZoneContextManager - supports asynchronous operations - optional
    contextManager: new ZoneContextManager(),
});

// Registering instrumentations (optional)
registerInstrumentations({
    instrumentations: [
        // This will include a number of events about the document load time, etc.
        // new DocumentLoadInstrumentation(),

        // Create traces for certain DOM events
        // new UserInteractionInstrumentation({
        //     eventNames: ['submit', 'click'],
        // }),

        // Creates traces for XHR requests
        // new XMLHttpRequestInstrumentation({
        //     propagateTraceHeaderCorsUrls: ['http://localhost:8000']
        // }),

        // Creates traces for Fetch API requests
        // new FetchInstrumentation()
    ],
});
