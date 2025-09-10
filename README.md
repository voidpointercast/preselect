# preselect

A Go application that scans CSV, JSON, and YAML files for user-defined keywords using string similarity metrics. The project is structured with a layered architecture consisting of API, Business, and Data layers.

This repository currently contains the basic skeleton of the application.

## Layers

- **API Layer**: Handles configuration and orchestrates application flow.
- **Business Layer**: Contains similarity checks, result processing, and a `DataSource` interface for pluggable inputs.
- **Data Layer**: Reads data from files in chunks for scalable processing.

## Usage

The project is in its early stages and presently only provides a scaffold for future development.

