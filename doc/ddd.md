
# Domain Driven Design

> Domain-Driven Design is an approach to software development that centers the development on programming a domain model that has a rich understanding of the processes and rules of a domain. The name comes from a 2003 book by Eric Evans that describes the approach through a catalog of patterns.

[Martin Fowler blog on domain drive design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

## Objectives
The aims of this document is to explain briefly the pattern and mostly the design choice found in this repository.

## Principles

Some choices in this repository can be odd at the first glance if not put in context, here is a list of view principles
choose to design the services:
* It is ok to write more code if it make it explicit
* Use layer-based architecture (with principles close to Clean Architecture, Ports and adapters, etc.)
    ```
         -------------------------->  [model]  <-----------------------
        |                                |                            |
    [transport] <--(direct call)--> [usecase] <-- (interfaces) --> [infra]
  
  ```
* Every layer should be the most independent to allow tests
* `domain` is an important part of the service, its value is there and should be the most explicit possible on the names

## Services Architecture

the hierarchy is design around modularity and single responsibility principles

- `main.go` => entry point of the service, this is where most of the dependency inject take place
        
- `infra/` =>
this folder is responsible to handle all the interaction outside the service on its most "obvious form". this mean only API call, database writting or reading, etc.
into package which describe the intent of it and the technology. like `postgres` and `etcd` are BAD PACKAGE name, but can be put in the file name inside a package like `users_repository/` or `cardpayer/`.
 the aims are to delegate technology choice to the very end and make this layer as "dumb" as possible to allow an easier manipulation such as add caching, tracing, or more advance pattern as data migration.
    * _You MUST NOT return object or value which leak the technology layer used underneath (such as SQL rows or things like this)._
    * _You MUST NOT import files from `domain/usecase`_
 
 
- `domain/` => this folder try to contain all the logic the service as been written for in the first place
    - `model/`
    this package contains pure data object which "model" the business object (or model) which will be used as "contract" and will be used by the `infra/` as data to send or save and in the `usecase` on data which will be used to run the business logic.
        * _You MUST NOT import from `domain/usecase` or `infra/`_
    - `usecase/`
    this package hold "features" packages, where each one is a complete set of "usecases" such as "send a notification to the user" or "retrieve preferences" 
        - `meterreadreminder/` an example of feature package name
 
- `transport` => this contains all the "transport" layer to communicate to the service, to send him commands separate by technology or "protocol"
    - `http`, `grpc`, `kafka` => are example of sub folder found in the `transport/` folder.
