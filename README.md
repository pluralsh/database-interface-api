# Database Interface API

This repository hosts the API defintion of the Custom Resource Definitions (CRD) used for the Database Interface project.
The provisioned unit of storage is a `Database`. The following CRDs are defined for managing the lifecycle of Databases:

- DatabaseRequest - Represents a request to provision a Database
- DatabaseClass - Represents a class of Datbase with similar characteristics
- Database - Represents a Database

The following CRDs are defined for managing the lifecycle of workloads accessing the Database:

- DatabaseAccessClass - Represents a class of accessors with similar access requirements
- DatabaseAccess - Represents an access secret to the Database