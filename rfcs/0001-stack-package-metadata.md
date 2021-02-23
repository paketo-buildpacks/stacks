# Add io.paketo.stack.packages label

## Proposal

We want to add a new label (io.paketo.stack.packages) to all stacks that contains a collection of metadata for all of the packages in the stack.

## Motivation

This metadata will help users have a more accurate depiction of exactly what is in the stack.

## Implementation


Schema:
```json
"io.paketo.stack.packages": "[{
                               "name": "<NAME>",
                               "version": "<VERSION>>",
                               "arch": "<ARCHITECTURE>>",
                               "summary": "<SUMMARY>",
                               "sourcePackage": {
                                 "name": "<NAME>",
                                 "version": "<VERSION>",
                                 "upstreamVersion": "<UPSTREAM_VERSION>"
                               }
                            }]"
```

Example:
```json
"io.paketo.stack.packages": "[{
                               "name": "libc6",
                               "version": "2.27-3ubuntu1.4",
                               "arch": "amd64",
                               "summary": "GNU C Library: Shared libraries",
                               "sourcePackage": {
                                 "name": "glibc",
                                 "version": "2.27-3ubuntu1.4",
                                 "upstreamVersion": "2.27"
                               }
                             }]"
```


## Unresolved Questions and Bikeshedding

Is there a clearer/more accurate label name we can use? Technically "stack" refers to a pair of images but we're only trying to represent the packages on a single image.


{{REMOVE THIS SECTION BEFORE RATIFICATION!}}
