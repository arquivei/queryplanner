# QueryPlanner

### A Golang library for planning queries in a structured order

---------------------

## Table of Contents

  - [1. Description](#Description)
  - [2. Technology Stack](#TechnologyStack)
  - [3. Getting Started](#GettingStarted)
  - [4. Changelog](#Changelog)
  - [5. Collaborators](#Collaborators)
  - [6. Contributing](#Contributing)
  - [7. Versioning](#Versioning)
  - [8. License](#License)
  - [9. Contact Information](#ContactInformation)

## <a name="Description" /> 1. Description

QueryPlanner is a generic library that aims to provide a framework for structuring queries that might need to reach different services or go through different steps to enrich a response. The library provides a way of describing the query's steps(IndexProviders/FieldProviders) and it's dependencies.

## <a name="TechnologyStack" /> 2. Technology Stack

| **Stack**     | **Version** |
|---------------|-------------|
| Golang        | v1.18       |
| golangci-lint | v1.46.2     |

## <a name="GettingStarted" /> 3. Getting Started

- ### <a name="Prerequisites" /> Prerequisites

  - Any [Golang](https://go.dev/doc/install) programming language version installed, preferred 1.18 or later.

- ### <a name="Install" /> Install
  
  ```
  go get -u github.com/arquivei/queryplanner
  ```

- ### <a name="ConfigurationSetup" /> Configuration Setup

  ```
  go mod vendor
  go mod tidy
  ```

- ### <a name="Usage" /> Usage
  
  - Import the package

    ```go
    import (
        "github.com/arquivei/queryplanner"
    )
    ```

- ### <a name="Examples" /> Examples
  
  - [Sample usage](https://github.com/arquivei/queryplanner/blob/master/examples/main.go)

## <a name="Changelog" /> 4. Changelog

  - **queryplanner 0.1.0 (May 31, 2022)**
  
    - [New] Documents: Code of Conduct, Contributing, License and Readme.
    - [New] Setting github's workflow with golangci-lint
    - [New] Decoupling this package from Arquivei's API projects.

## <a name="Collaborators" /> 5. Collaborators

- ### <a name="Authors" /> Authors

  <!-- markdownlint-disable -->
  <!-- prettier-ignore-start -->
	<table>
	<tr>
		<td align="center"><a href="https://github.com/victormn"><img src="https://avatars.githubusercontent.com/u/9757545?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Victor Nunes</b></sub></a></td>
	</tr>
	</table>
  <!-- markdownlint-restore -->
  <!-- prettier-ignore-end -->

- ### <a name="Maintainers" /> Maintainers
  
  <!-- markdownlint-disable -->
  <!-- prettier-ignore-start -->
	<table>
	<tr>
		<td align="center"><a href="https://github.com/marcosbmf"><img src="https://avatars.githubusercontent.com/u/34271729?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Marcos Barros</b></sub></a></td>
	</tr>
	</table>
  <!-- markdownlint-restore -->
  <!-- prettier-ignore-end -->

## <a name="Contributing" /> 6. Contributing

  Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## <a name="Versioning" /> 7. Versioning

  We use [Semantic Versioning](http://semver.org/) for versioning. For the versions
  available, see the [tags on this repository](https://github.com/arquivei/queryplanner/tags).

## <a name="License" /> 8. License
  
This project is licensed under the BSD 3-Clause - see the [LICENSE.md](LICENSE.md) file for details.

## <a name="ContactInformation" /> 9. Contact Information

  All contact may be doing by [marcos.filho@arquivei.com.br](mailto:marcos.filho@arquivei.com.br)