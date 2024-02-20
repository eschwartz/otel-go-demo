import {SpanStatusCode, trace} from "@opentelemetry/api";

let tracer = trace.getTracer('test-tracer');


// Wrapper around tracer.startActiveSpan()
// with added support for error handling and sending the span.
const withActiveSpan = (name, fn) => {
    return tracer.startActiveSpan(name, async span => {
        try {
            return await fn(span)
        }
        catch (err) {
            span.setStatus({ code: SpanStatusCode.ERROR, message: err.message });
            span.setAttributes({
                "app.error": err.message,
                "app.error.stack": err.stack,
            })
            throw err
        }
        finally {
            span.end()
        }
    })
}

// Request handler, instrumented with otel
// adapted from https://github.com/honeycombio/example-greeting-service/blob/main/web/src/index.js
const request = async (url, opts = {}) => {
    let method = opts.method || 'GET';

    return withActiveSpan(`Request: ${method} ${url}`, async span => {
        span.setAttributes({
            'request.method': method,
            'request.url': url,
        })
        const res = await fetch(url, {
            ...opts,
            // Add traceparent header for trace propagation
            headers: {
                ...(opts.headers || {}),
                traceparent: `00-${span.spanContext().traceId}-${span.spanContext().spanId}-01`,
            }
        })

        span.setAttributes({
            'response.status_code': res.status,
        })

        return res
    })
}




document.getElementById("search").onsubmit = async (evt) => {
    evt.preventDefault()
    withActiveSpan("submit fetch", async span => {
        const limit = document.getElementById("limit").value
        const term = document.getElementById("term").value

        span.setAttributes({
            "app.fetchTrigger": "search",
            "app.itemsSearch.limit": limit,
            "app.itemsSearch.term": term,
        })

        const res = await request(`/api/items?limit=${limit}&q=${term}`)
        const items = await res.json()

        span.setAttributes({
            'app.resultCount': items.length,
            'app.resultJson': JSON.stringify(items)
        })

        renderItems(items)
    })

    return false;
}

document.getElementById("lucky").onclick = async () => {
    withActiveSpan("submit fetch", async span => {
        // Pick a random limit between 0 and 3
        const randomLimit = Math.round(Math.random() * 3)
        // bicycle is a sufficiently random enough term, see https://xkcd.com/221/
        const randomTerm = "bicycle";

        span.setAttributes({
            "app.fetchTrigger": "search",
            "app.itemsSearch.limit": randomLimit,
            "app.itemsSearch.term": randomTerm,
        })
        const res = await request(`/api/items?limit=${randomLimit}&q=${randomTerm}`)
        const items = await res.json()

        span.setAttributes({
            'app.resultCount': items.length,
            'app.resultJson': JSON.stringify(items)
        })

        renderItems(items)
    })
}



function renderItems(items) {
    document.getElementById("items")
        .innerHTML = items.map(item => `<li>${item.value}</li>`).join('')
}
