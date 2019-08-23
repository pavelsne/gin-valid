# Contribution guide: Adding new validator

The following is a list of everything needed to add a new validator to the service.  The placeholder name `V` should be replaced with the name of the validator in the example function and variable names.  More detailed descriptions of each requirement can be found in the sections below.
- A validation function, `validateV()`, in the `internal/web/validate.go` file.
- A `v_results.go` file that contains a template to render the results of the validation.  The template should be stored in a const string called `VResults`.
- A `renderVResults()` function, in the `internal/web/results.go` file, that uses the `VResults` template to render the results.
- Configuration settings for the new validator:
    - `ServerCfg.Executables.V` should point to the executable that runs the validation.
    - `ServerCfg.Settings.Validators` should include the (all lowercase) name of the validator.
- The executable should be included in the Dockerfile.


## Validation function

Validation functions are named `validateV()`, where `V` is the name of the validator.  This function should take two arguments and return an `error` type.
The two arguments are:
- `valroot`: The location of the repository that will be validated.  Use this directory as the starting point for the validation command.
- `resdir`: The results directory where the results of the validation should be stored.  The validator function should create two files in this directory: a badge and a file with the results.

For the name of the badge file name, use `srvcfg.Label.ResultsBadge`.  Depending on the results of the validation, the contents of this file should be one of the const strings found in `internal/resources/svg.go`.

The format of the results is different for each validator.  These results will be processed by the `VResults` function to render the `v_results.go` template, so the results should be stored in a way that will make this most convenient.  The name of the file should be `srvcfg.Label.ResultsFile`.

Once the validation function has been written, a `switch` case should be added for it at the bottom of the `runValidator()` function.


## Results template

The template should contain a header with the badge and name of the repository.  The main body should be the rendered contents of the results.

See the existing templates for examples on what this should look like.

The name of this template should be `VResults`.


## Results rendering function

The `renderVResults()` function should use the data stored in the results file (`resdir/srvcfg.Label.ResultsFile`) to render the results page.  This function should take 6 arguments.
The arguments are:
- `w` and `r`: The `http.ResponseWriter` and `http.Request` coming from the web request.  Use these to render the resulting page.
- `badge`: A byte slice containing the badge contents.  The template should use the data in this slice to render the badge.
- `content`: A byte slice containing the contents of the results file that you stored in the [Validation function](#validation-function).
- `user` and `repo`: The user and repository names as strings.  Use these to render the repository name in the header and for error reporting.

This function should load the main layout template (found in `templates.Layout`), then parse the validator template (found in `templates.VResults`), add the data it requires and execute it.


## Configuration settings

These settings can be added at runtime, but they can also be added to the default configuration for simplicity.


## Dockerfile

The executable and any required dependencies should be included in the `RUNNER IMAGE` part of the Dockerfile.  Like with NIX, binaries that are built from source can be built in separate images then copied to the main Docker runner image.  See the `NIX BUILDER IMAGE` section as well as the [Multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/) Docker documentation.
