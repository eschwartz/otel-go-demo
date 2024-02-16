import {trace} from "@opentelemetry/api";

trace.getTracer('test-tracer')
    .startActiveSpan('test span client', span => {
        setTimeout(() => {
            console.log('done')
            span.end()
        }, 1000)
    })
