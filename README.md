# VWAP calculation engine

## Running the project

This project has a `Makefile` on the root folder that has the following commands:

- `make test`: Run all tests and generate a coverage file
- `make run`: Start the application ~~(use docker stop to stop the execution)~~

## Development environment

The development environement used is based on VSCode Dev Containers, that already provide a series of pre-configured details to run the project, lint, tests, and debug.

## TODO

- [ ] Handle better with the float vars
- [ ] Create concise services instead to glue all on main func
- [ ] Improve the error handling
- [ ] Improve logs
- [X] Fix issue using the Makefile that not allow to stop the application
