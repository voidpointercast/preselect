# preselect

A Go application that scans CSV, JSON, and YAML files for user-defined keywords using string similarity metrics. The project is structured with a layered architecture consisting of API, Business, and Data layers.

This repository currently contains the basic skeleton of the application.

## Layers

- **API Layer**: Handles configuration and orchestrates application flow.
- **Business Layer**: Contains similarity checks, result processing, and a `DataSource` interface for pluggable inputs.
- **Data Layer**: Reads data from files in chunks for scalable processing.

## Usage

`preselect` scans a directory for files whose extensions are associated with a data
source. Mappings are supplied on the command line using the `-ext` flag:

```
preselect -dir ./data -ext txt=,;
```

The above command scans `./data` and tokenizes `.txt` files using commas and
semicolons as delimiters. When no mapping is provided for `txt` files, they are
tokenized using spaces and newlines by default.

The project is in its early stages and presently only provides a scaffold for
future development.

