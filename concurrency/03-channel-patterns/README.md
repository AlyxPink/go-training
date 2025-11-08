# Exercise 03: Channel Patterns

**Difficulty**: ⭐⭐⭐
**Estimated Time**: 60 minutes

## Learning Objectives
- Implement pipeline patterns
- Master fan-out/fan-in
- Build worker pools with channels
- Understand generator patterns

## Problem Description
Build common concurrent patterns using channels: pipelines for data transformation, fan-out for parallel processing, and fan-in for result aggregation.

## Tasks
1. Pipeline: chain processing stages
2. Fan-out: distribute work to multiple workers
3. Fan-in: aggregate results from multiple sources
4. Generator: produce values on-demand

## Concepts Covered
- Pipeline composition
- Fan-out/fan-in patterns
- Channel-based worker pools
- Generator functions

## Testing
```bash
go test -race -v
```
