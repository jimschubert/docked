/*
Package docked provides types and functionality for analyzing and linting Dockerfiles.

docked uses the Docker buildkit parser to retrieve the AST of an input Dockerfile.
It also provides a simple API for defining and registering rules for processing of the AST.
All in-built rules are built upon this API. See those defined under the validations package.

Configuration

An external YAML configuration is supported by docked.Config. The configuration allows for ignoring in-built
rules, overriding priority of in-built rules, as well as defining custom rules based on the
validations.SimpleRegexRule structure.

Analysis

Invoking docked.Docked#Analysis will use the list of in-built validation rules, and return a docked.AnalysisResult.
The result should be walked programmatically to generate a report. Please see reports under the reporting package for examples.
The HTML and JSON reporters under the reporter package provide implementations for use in the
accompanying cli tool for use in CI/CD pipelines.
*/
package docked
