# TinyString Code Reduction and Refactoring Prompt

## Objective
Reduce code lines and eliminate unnecessary code while maximizing reusability through small functions. Priority is code reduction and architecture improvement for better maintainability.

## Key Requirements

1. **Code Reduction Focus**: Eliminate unnecessary lines of code
2. **Maximum Reusability**: Use small, reusable functions  
3. **Minimal Private Elements**: Reduce private methods and variables to minimum
4. **Architecture Improvement**: Simplicity, fewer lines of code, and good performance
5. **Public API Stability**: Public methods must not change, only private ones
6. **Error Handling Consolidation**: Keep error handling in `error.go`, improve ErrorF and Fmt/sprintf design

## Current Problems Identified
- Redundant and repetitive code in Fmt/sprintf functions
- Poor error handling design between ErrorF and Fmt
- Repetitive type switches that can be replaced with generics
- Excessive private methods and helper functions

## Success Metrics
- Reduce total lines of code by 30-50%
- Maintain or improve WebAssembly binary size reduction
- Simplify architecture for easier maintenance
- Maintain all existing public API functionality

## Strategy Reference
See `ISSUE_GENERIC_ARQ.md` for detailed implementation strategy and progress tracking.
