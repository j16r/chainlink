# WASM Chainlink

This is an experinent to build a version of chainlink that only supports pure WASM jobs.

JobSpecs currently look a bit like this:

    {
      "initiators": [{"type": "web"}],
      "tasks": [
        {"type": "httpget", "params": {"get": "https://bitstamp.net/api/ticker/"}},
        {"type": "jsonparse", "params": {"path": ["last"]}},
        {
          "type": "ethtx",
          "confirmations": 0,
          "params": {
            "address": "0xaa664fa2fdc390c662de1dbacf1218ac6e066ae6",
            "functionSelector": "setBytes(bytes32,bytes)",
            "format": "bytes",
            "dataPrefix": "0x0000000000000000000000000000000000000000000000000000000000000001"
          }
        }
      ]
    }

A WASM job spec would look a bit more like this:

    {
      "parameters": {
        "url": {
          "type": "string",
          "default": "https://bitstamp.net/api/ticker/"
        }
        "contractAddress": {
          "type": "string",
          "default": "0xaa664fa2fdc390c662de1dbacf1218ac6e066ae6"
        }
      }
      "wasm": "<base64 encoded wasm binary>"
    }

This version of chainlink would not support any tasks or adapters, it would
totally be up to a user to define their entire job as a wasm program. However
the WASM program would need to reach out to the outside world, so we'd provide
a bunch of function globals that can be called from your wasm program. These
would look a lot like the adapters.

## Parameters

Parameters can be overwritten by the job initiator,,

## Functions

  * sendTransaction(address, value, data)
  * httpGet(address)
  * httpPost(address)
  * httpPut(address)

In theory we'd need fewer functions because our WASM code can be turing
complete, things like parsing JSON, arithmetic operations, *sleeping etc. Can
all be implemented natively in wasm.

* Sleep may require a function to be implemented efficiently.

## Returning

It's easy for a wasm program to return some arbitrary data, this might allow
for a nice optimization. Whatever we return from the wasm program could
automatically be written to the `fulfill` function of the specified contract.
Rather than have contract writers explictly add this step.

## Accounting

Since WASM programs are turing complete, we face the halting problem. A WASM
program could theoretically have an infinite loop. So we'll attach some kind of
cost to instructions and terminate execution if the payment for the job is
exhausted.

## SGX

This implementation also makes it easier to protect our computation in SGX,
because the entire off chain operation can be peformed from within an SGX
enclave.

# Initiating WASM jobs

One advantage of doing WASM jobs is that we could make initiation really
flexible, I propose that rather than setting an initiator in the job spec. You
handle initiation explicitly through calls into chainlink. e.g:

    httpListen("suffix")

httpListen would effectively block the WASM Task in an efficient way and
"mount" a handler under chainlink's own HTTP handler. e.g:

    /v2/jobs/${guid}/suffix

${guid} = id of the job spec returned from chainlink.

There are numerous other examples that would work:

    logListen
    RPCListen
    WSListen
    EmailListen

Or if we wanted to know more about the initiaton in advance, we could inspect
the symbols of the WASM and look for some that match a certain name / arity
that indicate they are "initiator entry points".

# Storing State

One useful thing to add might be to allow WASM jobs to save state. Any return
value from a WASM job could be provided as an argument to any further
invocations of a WASM adapter. We would likely have a hard limit on the amount
of state that can be stored, say 4k per job. It would be formatted as JSON.

# Pros/Cons

What are the tradeoffs between the WASM chainlink and the the job spec
descriptions that we currently use?

Cons:

  * Turing complete
  * WASM is more opaque.
  * Less progress can be communicated back to the user. We don't really know
    what the "path" of a WASM binary is supposed to look like so we can't tell
    if you're 100% or 0%. Only guess.
  * A failing WASM task might have to present a stack trace or some other
    output, which would be much harder to debug.
  * More opportunity for "bombs", e.g: WASM payloads that are harmful, use lots
    of resources or otherwise.
  * More complexity needed to account for cost of executing WASM.

Pros:

  * Turing complete.
  * WASM should be faster, even when interpereted.
  * Use any language you prefer, easier to port existing code based to smart
    contract model.
  * More of the job can occur in SGX.
  * WASM jobs can be more composeable.
  * Jobs are much more flexible, can have logic, loops, multiple transactions,
    no transaction etc.
  * This design makes external inititators and bridges largely obsolete.
