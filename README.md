# test_task_16
The some test task

## Description
Write some code, that will flatten an array of arbitrarily nested arrays of integers into a flat array of integers. e.g. [[1,2,[3]],4] -> [1,2,3,4].

Your solution should be a link to a gist on gist.github.com with your implementation.

When writing this code, you can use any language you're comfortable with. The code must be well tested and documented. Please include unit tests and any documentation you feel is necessary. In general, treat the quality of the code as if it was ready to ship to production.

Try to avoid using language defined methods like Ruby's Array#flatten.

## Implementation
```
pkg/flatter
```

## Test
```
go test -v pkg/flatter/*.go
```
