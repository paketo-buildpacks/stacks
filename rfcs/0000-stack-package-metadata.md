# Add io.paketo.stack.packages label

## Proposal

We want to add a new label (io.paketo.stack.packages) to all stacks that contains a collection of metadata for all of the packages in the stack.

## Motivation

This metadata will help users have a more accurate depiction of exactly what is in the stack.

## Implementation


Schema:
```json
"io.paketo.stack.packages": "[{"name":"<NAME>","version":"<VERSION>>","arch":"<ARCHITECTURE>>","description":"<DESCRIPTION>>"}]"
```

Example:
```json
"io.paketo.stack.packages": "[{"name":"base-files","version":"10.1ubuntu2.10","arch":"amd64","description":"Secure Sockets Layer toolkit - cryptographic utility"}]"
```


## Unresolved Questions and Bikeshedding

Is there a clearer/more accurate label name we can use? Technically "stack" refers to a pair of images but we're only trying to represent the packages on a single image.


{{REMOVE THIS SECTION BEFORE RATIFICATION!}}
