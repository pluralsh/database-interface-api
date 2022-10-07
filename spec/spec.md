# Database Interface

Authors:

* Lukasz Zajaczkowski [@zreigz](https://github.com/zreigz)


## Notational Conventions

The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119](http://tools.ietf.org/html/rfc2119) (Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels", BCP 14, RFC 2119, March 1997).

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard](http://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18).

An implementation is not compliant if it fails to satisfy one or more of the MUST, REQUIRED, or SHALL requirements for the protocols it implements.
An implementation is compliant if it satisfies all the MUST, REQUIRED, and SHALL requirements for the protocols it implements.

## Objective

To define a standard that enables a database vendor to develop an RPC-based plugin once and have it work across multiple Kubernetes clusters.


### Database Lifecycle

```
     CreateDatabase +-------------+ DeleteDatabase 
    +------------->|  CREATED   +--------------+ 
    |              +---+----^---+              |    
    |          Grant   |    | Revoke           v    
    +++        Database|    | Database        +++  
    |X|        Access  |    | Access          | |  
    +-+            +---v----+---+             +-+  
                   |   BOUND    |                  
                   +---+----^---+                  

```


## Database Interface

This section describes the interface between Database systems and Plugins.

### RPC Interface


```protobuf
syntax = "proto3";
package database.v1alpha1;

import "google/protobuf/descriptor.proto";

option go_package = "pluralsh/database-interface-spec;database";

extend google.protobuf.EnumOptions {
    // Indicates that this enum is OPTIONAL and part of an experimental
    // API that may be deprecated and eventually removed between minor
    // releases.
    bool alpha_enum = 1116;
}

extend google.protobuf.EnumValueOptions {
    // Indicates that this enum value is OPTIONAL and part of an
    // experimental API that may be deprecated and eventually removed
    // between minor releases.
    bool alpha_enum_value = 1116;
}

extend google.protobuf.FieldOptions {
    // Indicates that a field MAY contain information that is sensitive
    // and MUST be treated as such (e.g. not logged).
    bool database_secret = 1115;

    // Indicates that this field is OPTIONAL and part of an experimental
    // API that may be deprecated and eventually removed between minor
    // releases.
    bool alpha_field = 1116;
}

extend google.protobuf.MessageOptions {
    // Indicates that this message is OPTIONAL and part of an experimental
    // API that may be deprecated and eventually removed between minor
    // releases.
    bool alpha_message = 1116;
}

extend google.protobuf.MethodOptions {
    // Indicates that this method is OPTIONAL and part of an experimental
    // API that may be deprecated and eventually removed between minor
    // releases.
    bool alpha_method = 1116;
}

extend google.protobuf.ServiceOptions {
    // Indicates that this service is OPTIONAL and part of an experimental
    // API that may be deprecated and eventually removed between minor
    // releases.
    bool alpha_service = 1116;
}

service Identity {
    // This call is meant to retrieve the unique provisioner Identity.
    // This identity will have to be set in DatabaseClaim.DriverName field in order to invoke this specific provisioner.
    rpc DriverGetInfo (DriverGetInfoRequest) returns (DriverGetInfoResponse) {}
}

service Provisioner {
    // This call is made to create the database in the backend.
    // This call is idempotent
    //    1. If a database that matches both name and parameters already exists, then OK (success) must be returned.
    //    2. If a database by same name, but different parameters is provided, then the appropriate error code ALREADY_EXISTS must be returned.
    rpc DriverCreateDatabase (DriverCreateDatabaseRequest) returns (DriverCreateDatabaseResponse) {}
    // This call is made to delete the database in the backend.
    // If the database has already been deleted, then no error should be returned.
    rpc DriverDeleteDatabase (DriverDeleteDatabaseRequest) returns (DriverDeleteDatabaseResponse) {}

    // This call grants access to an account. The account_name in the request shall be used as a unique identifier to create credentials.
    // The account_id returned in the response will be used as the unique identifier for deleting this access when calling DriverRevokeDatabaseAccess.
    rpc DriverGrantDatabaseAccess (DriverGrantDatabaseAccessRequest) returns (DriverGrantDatabaseAccessResponse);
    // This call revokes all access to a particular database from a principal.
    rpc DriverRevokeDatabaseAccess (DriverRevokeDatabaseAccessRequest) returns (DriverRevokeDatabaseAccessResponse);
}

// S3SignatureVersion is the version of the signing algorithm for all s3 requests
enum S3SignatureVersion {
    UnknownSignature = 0;
    // S3V2, Signature version v2
    S3V2 = 1;
    // S3V4, Signature version v4
    S3V4 = 2;
}

enum AnonymousDatabaseAccessMode {
    UnknownDatabaseAccessMode = 0;
    // Default, disallow uncredentialed access to the backend storage.
    Private = 1;
    // Read only, uncredentialed users can call ListDatabase and GetObject.
    ReadOnly = 2;
    // Write only, uncredentialed users can only call PutObject.
    WriteOnly = 3;
    // Read/Write, uncredentialed users can read objects as well as PutObject.
    ReadWrite = 4;
}

enum AuthenticationType {
    UnknownAuthenticationType = 0;
    // Default, KEY based authentication.
    Key = 1;
    // Storageaccount based authentication.
    IAM = 2;
}

message S3 {
    // region denotes the geographical region where the S3 server is running
    string region = 1;
    // signature_version denotes the signature version for signing all s3 requests
    S3SignatureVersion signature_version = 2;
}

message AzureBlob {
    // storage_account is the id of the azure storage account
    string storage_account = 1;
}

message GCS {
    // private_key_name denotes the name of the private key in the storage backend
    string private_key_name = 1;
    // project_id denotes the name of the project id in the storage backend
    string project_id = 2;
    // service_account denotes the name of the service account in the storage backend
    string service_account = 3;
}

message Protocol {
    oneof type {
        S3 s3 = 1;
        AzureBlob azureBlob = 2;
        GCS gcs = 3;
    }
}

message CredentialDetails {
    // map of the details in the secrets for the protocol string
    map<string, string> secrets = 1;
}

message DriverGetInfoRequest {
    // Intentionally left blank
}

message DriverGetInfoResponse {
    // This field is REQUIRED
    // The name MUST follow domain name notation format
    // (https://tools.ietf.org/html/rfc1035#section-2.3.1). It SHOULD
    // include the plugin's host company name and the plugin name,
    // to minimize the possibility of collisions. It MUST be 63
    // characters or less, beginning and ending with an alphanumeric
    // character ([a-z0-9A-Z]) with dashes (-), dots (.), and
    // alphanumerics between.
    string name = 1;
}

message DriverCreateDatabaseRequest {
    // This field is REQUIRED
    // name specifies the name of the database that should be created.
    string name = 1;

    // This field is OPTIONAL
    // The caller should treat the values in parameters as opaque. 
    // The receiver is responsible for parsing and validating the values.
    map<string,string> parameters = 2;
}

message DriverCreateDatabaseResponse {
    // database_id returned here is expected to be the globally unique 
    // identifier for the database in the object storage provider.
    string database_id = 1;

    // database_info returned here stores the data specific to the
    // database required by the object storage provider to connect to the database.
    Protocol database_info = 2;
}

message DriverDeleteDatabaseRequest {
    // This field is REQUIRED
    // database_id is a globally unique identifier for the database
    // in the object storage provider 
    string database_id = 1;
}

message DriverDeleteDatabaseResponse {
    // Intentionally left blank
}

message DriverGrantDatabaseAccessRequest {
    // This field is REQUIRED
    // database_id is a globally unique identifier for the database
    // in the object storage provider 
    string database_id = 1;

    // This field is REQUIRED
    // name field is used to define the name of the database access object.
    string name = 2;

    // This field is REQUIRED
    // Requested authentication type for the database access.
    // Supported authentication types are KEY or IAM.
    AuthenticationType authentication_type = 3;

    // This field is OPTIONAL
    // The caller should treat the values in parameters as opaque. 
    // The receiver is responsible for parsing and validating the values.
    map<string,string> parameters = 4;
}

message DriverGrantDatabaseAccessResponse {
    // This field is REQUIRED
    // This is the account_id that is being provided access. This will
    // be required later to revoke access. 
    string account_id = 1;

    // This field is REQUIRED
    // Credentials supplied for accessing the database ex: aws access key id and secret, etc.
    map<string, CredentialDetails> credentials = 2;
}

message DriverRevokeDatabaseAccessRequest {
    // This field is REQUIRED
    // database_id is a globally unique identifier for the database
    // in the object storage provider.
    string database_id = 1;

    // This field is REQUIRED
    // This is the account_id that is having its access revoked.
    string account_id = 2;
}

message DriverRevokeDatabaseAccessResponse {
    // Intentionally left blank
}

```

##### Size Limits

The general size limit for a particular field MAY be overridden by specifying a different size limit in said field's description.
Unless otherwise specified, fields SHALL NOT exceed the limits documented here.
These limits apply for messages generated by both Database systems and plugins.

| Size       | Field Type          |
|------------|---------------------|
| 128 bytes  | string              |
| 4 KiB      | map<string, string> |

##### `REQUIRED` vs. `OPTIONAL`

* A field noted as `REQUIRED` MUST be specified, subject to any per-RPC caveats; caveats SHOULD be rare.
* A `repeated` or `map` field listed as `REQUIRED` MUST contain at least 1 element.
* A field noted as `OPTIONAL` MAY be specified and the specification SHALL clearly define expected behavior for the default, zero-value of such fields.

Scalar fields, even REQUIRED ones, will be defaulted if not specified and any field set to the default value will not be serialized over the wire as per [proto3](https://developers.google.com/protocol-buffers/docs/proto3#default).

#### Timeouts

Any of the RPCs defined in this spec MAY timeout and MAY be retried.
The Database system MAY choose the maximum time it is willing to wait for a call, how long it waits between retries, and how many time it retries (these values are not negotiated between plugin and Database system).

Idempotency requirements ensure that a retried call with the same fields continues where it left off when retried.
The only way to cancel a call is to issue a "negation" call if one exists.
For example, issue a `DeleteDatabase` call to cancel a pending `CreateDatabase` operation, etc.

### Error Scheme

All Database API calls defined in this spec MUST return a [standard gRPC status](https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto).
Most gRPC libraries provide helper methods to set and read the status fields.

The status `code` MUST contain a [canonical error code](https://github.com/grpc/grpc-go/blob/master/codes/codes.go). Database systems MUST handle all valid error codes. Each RPC defines a set of gRPC error codes that MUST be returned by the plugin when specified conditions are encountered. In addition to those, if the conditions defined below are encountered, the plugin MUST return the associated gRPC error code.

| Condition | gRPC Code | Description | Recovery Behavior |
|-----------|-----------|-------------|-------------------|
| Missing required field | 3 MISSING_ARGUMENT | Indicates that a required field is missing from the request. More human-readable information MAY be provided in the `status.message` field. | Caller MUST fix the request by adding the missing required field before retrying. |
| Invalid or unsupported field in the request | 3 INVALID_ARGUMENT | Indicates that the one or more fields in this field is either not allowed by the Plugin or has an invalid value. More human-readable information MAY be provided in the gRPC `status.message` field. | Caller MUST fix the field before retrying. |
| Permission denied | 7 PERMISSION_DENIED | The Plugin is able to derive or otherwise infer an identity from the secrets present within an RPC, but that identity does not have permission to invoke the RPC. | System administrator SHOULD ensure that requisite permissions are granted, after which point the caller MAY retry the attempted RPC. |
| Operation pending for database | 10 ABORTED | Indicates that there is already an operation pending for the specified database. In general the Cluster Orchestrator (CO) is responsible for ensuring that there is no more than one call "in-flight" per database at a given time. However, in some circumstances, the Database system MAY lose state (for example when the Database system crashes and restarts), and MAY issue multiple calls simultaneously for the same database. The Plugin, SHOULD handle this as gracefully as possible, and MAY return this error code to reject secondary calls. | Caller SHOULD ensure that there are no other calls pending for the specified volume, and then retry with exponential back off. |
| Call not implemented | 12 UNIMPLEMENTED | The invoked RPC is not implemented by the Plugin or disabled in the Plugin's current mode of operation. | Caller MUST NOT retry. Caller MAY call `GetPluginInfo` to discover Plugin info. |
| Not authenticated | 16 UNAUTHENTICATED | The invoked RPC does not carry secrets that are valid for authentication. | Caller SHALL either fix the secrets provided in the RPC, or otherwise regalvanize said secrets such that they will pass authentication by the Plugin for the attempted RPC, after which point the caller MAY retry the attempted RPC. |

The status `message` MUST contain a human readable description of error, if the status `code` is not `OK`.
This string MAY be surfaced by Database system to end users.

The status `details` MUST be empty. In the future, this spec MAY require `details` to return a machine-parsable protobuf message if the status `code` is not `OK` to enable Database system's to implement smarter error handling and fault resolution.

### Provisioner Service RPC

Provisioner service RPCs allow a Database system to query a plugin for information, create and delete database as well as grant and revoke access to the database to various principals running the workload.
The general flow of the success case MAY be as follows (protos illustrated in YAML for brevity):

1. Database system queries metadata via Identity RPC.

```
   # Database system --(DriverGetInfo)--> Plugin
   request:
   response:
      name: org.foo.whizbang.super-plugin
```
```
message DriverGetInfoRequest {
    // Intentionally left blank
}

message DriverGetInfoResponse {
    string name = 1;
}
```

#### `DriverCreateDatabase`

A Controller Plugin MUST implement this RPC call.
This RPC will be called by the Database system to provision a new database on behalf of a Database user.

This operation MUST be idempotent.
If a volume corresponding to the specified database `name` already exists, is accessible from `accessibility_requirements, and is compatible with the specified attributes of the database in the `DriverCreateDatabase`, the Plugin MUST reply `0 OK` with the corresponding `DriverCreateDatabaseResponse`.

```
message DriverCreateDatabaseRequest {  
    // Idempotency - This name is generated by the Database system to achieve
    // idempotency.  
    // This field is REQUIRED.
    string database_name = 1;

    // This field is OPTIONAL
    // Protocol specific information required by the call is passed in as key,value pairs.
    map<string,string> database_context = 2;
}

message DriverCreateDatabaseResponse {
    // Intentionally left blank
}

```
##### CreateDatabase Errors

If the plugin is unable to complete the CreateDatabase call successfully, it MUST return a non-ok gRPC code in the gRPC status.
If the conditions defined below are encountered, the plugin MUST return the specified gRPC error code.
The Database system MUST implement the specified error recovery behavior when it encounters the gRPC error code.

| Condition | gRPC Code | Description | Recovery Behavior |
|-----------|-----------|-------------|-------------------|
| Database Name invalid | 3 INVALID_ARGUMENT | Besides the general cases, this code MUST also be used to indicate when plugin cannot create a database due to parameter check on the backend. More human-readable information SHOULD be provided in the gRPC `status.message` field for the problem. | Caller should retry call with corrected parameter values.|
| Database already exists but is incompatible | 6 ALREADY_EXISTS | Indicates that a database corresponding to the specified database `name` already exists but is incompatible with the specified attributes. | Caller MUST fix the arguments or use a different `name` before retrying. |
| Unsupported value | 11 OUT_OF_RANGE | Indicates that the value requested cannot be used to provision the database. | Caller MUST fix the attributes before retrying. |


### Secrets Requirements (where applicable)

Secrets MAY be required by plugin to complete a RPC request.
A secret is a string to string map where the key identifies the name of the secret (e.g. "username" or "password"), and the value contains the secret data (e.g. "bob" or "abc123").
Each key MUST consist of alphanumeric characters, '-', '_' or '.'.
Each value MUST contain a valid string.
An SP MAY choose to accept binary (non-string) data by using a binary-to-text encoding scheme, like base64.
An SP SHALL advertise the requirements for required secret keys and values in documentation.
Database system SHALL permit passing through the required secrets.
A Database system MAY pass the same secrets to all RPCs, therefore the keys for all unique secrets that an SP expects MUST be unique across all Database operations.
This information is sensitive and MUST be treated as such (not logged, etc.) by the Database system.


## Protocol

### Connectivity

* A Database system SHALL communicate with a Plugin using gRPC to access the `Provisioner` service.
  * proto3 SHOULD be used with gRPC, as per the [official recommendations](http://www.grpc.io/docs/guides/#protocol-buffer-versions).
  * All Plugins SHALL implement the REQUIRED Identity service RPCs.
* The Database system SHALL provide the listen-address for the Plugin by way of the `Database_ENDPOINT` environment variable.
  Plugin components SHALL create, bind, and listen for RPCs on the specified listen address.
  * Only UNIX Domain Sockets MAY be used as endpoints.
    This will likely change in a future version of this specification to support non-UNIX platforms.
* All supported RPC services MUST be available at the listen address of the Plugin.

### Security

* The Database system operator and Provisioner Sidecar SHOULD take steps to ensure that any and all communication between the Database system and Plugin Service are secured according to best practices.
* Communication between a Database system and a Plugin SHALL be transported over UNIX Domain Sockets.
  * gRPC is compatible with UNIX Domain Sockets; it is the responsibility of the Database system operator and Provisioner Sidecar to properly secure access to the Domain Socket using OS filesystem ACLs and/or other OS-specific security context tooling.
  * SP’s supplying stand-alone Plugin controller appliances, or other remote components that are incompatible with UNIX Domain Sockets MUST provide a software component that proxies communication between a UNIX Domain Socket and the remote component(s).
    Proxy components transporting communication over IP networks SHALL be responsible for securing communications over such networks.
* Both the Database system and Plugin SHOULD avoid accidental leakage of sensitive information (such as redacting such information from log files).

### Debugging

* Debugging and tracing are supported by external, Database-independent additions and extensions to gRPC APIs, such as [OpenTracing](https://github.com/grpc-ecosystem/grpc-opentracing).

## Configuration and Operation

### General Configuration

* The `Database_ENDPOINT` environment variable SHALL be supplied to the Plugin by the Provisioner Sidecar.
* An operator SHALL configure the Database system to connect to the Plugin via the listen address identified by `Database_ENDPOINT` variable.
* With exception to sensitive data, Plugin configuration SHOULD be specified by environment variables, whenever possible, instead of by command line flags or bind-mounted/injected files.


#### Plugin Bootstrap Example

* Sidecar -> Plugin: `Database_ENDPOINT=unix:///path/to/unix/domain/socket.sock`.
* Operator -> Database system: use plugin at endpoint `unix:///path/to/unix/domain/socket.sock`.
* Database system: monitor `/path/to/unix/domain/socket.sock`.
* Plugin: read `Database_ENDPOINT`, create UNIX socket at specified path, bind and listen.
* Database system: observe that socket now exists, establish connection.
* Database system: invoke `GetPluginInfo.

#### Filesystem

* Plugins SHALL NOT specify requirements that include or otherwise reference directories and/or files on the root filesystem of the Database system.
* Plugins SHALL NOT create additional files or directories adjacent to the UNIX socket specified by `Database_ENDPOINT`; violations of this requirement constitute "abuse".
  * The Provisioner Sidecar is the ultimate authority of the directory in which the UNIX socket endpoint is created and MAY enforce policies to prevent and/or mitigate abuse of the directory by Plugins.

### Supervised Lifecycle Management

* For Plugins packaged in software form:
  * Plugin Packages SHOULD use a well-documented container image format (e.g., Docker, OCI).
  * The chosen package image format MAY expose configurable Plugin properties as environment variables, unless otherwise indicated in the section below.
    Variables so exposed SHOULD be assigned default values in the image manifest.
  * A Provisioner Sidecar MAY programmatically evaluate or otherwise scan a Plugin Package’s image manifest in order to discover configurable environment variables.
  * A Plugin SHALL NOT assume that an operator or Provisioner Sidecar will scan an image manifest for environment variables.

#### Environment Variables

* Variables defined by this specification SHALL be identifiable by their `Database_` name prefix.
* Configuration properties not defined by the Database specification SHALL NOT use the same `Database_` name prefix; this prefix is reserved for common configuration properties defined by the Database specification.
* The Provisioner Sidecar SHOULD supply all RECOMMENDED Database environment variables to a Plugin.
* The Provisioner Sidecar SHALL supply all REQUIRED Database environment variables to a Plugin.

##### `Database_ENDPOINT`

Network endpoint at which a Plugin SHALL host Database RPC services. The general format is:

    {scheme}://{authority}{endpoint}

The following address types SHALL be supported by Plugins:

    unix:///path/to/unix/socket.sock

Note: All UNIX endpoints SHALL end with `.sock`. See [gRPC Name Resolution](https://github.com/grpc/grpc/blob/master/doc/naming.md).

This variable is REQUIRED.

#### Operational Recommendations

The Provisioner Sidecar expects that a Plugin SHALL act as a long-running service vs. an on-demand, CLI-driven process.

Supervised plugins MAY be isolated and/or resource-bounded.

##### Logging

* Plugins SHOULD generate log messages to ONLY standard output and/or standard error.
  * In this case the Provisioner Sidecar SHALL assume responsibility for all log lifecycle management.
* Plugin implementations that deviate from the above recommendation SHALL clearly and unambiguously document the following:
  * Logging configuration flags and/or variables, including working sample configurations.
  * Default log destination(s) (where do the logs go if no configuration is specified?)
  * Log lifecycle management ownership and related guidance (size limits, rate limits, rolling, archiving, expunging, etc.) applicable to the logging mechanism embedded within the Plugin.
* Plugins SHOULD NOT write potentially sensitive data to logs (e.g. secrets).

##### Available Services

* Plugin Packages MAY support all or a subset of Database services; service combinations MAY be configurable at runtime by the Provisioner Sidecar.
  * A plugin MUST know the "mode" in which it is operating (e.g. node, controller, or both).
  * This specification does not dictate the mechanism by which mode of operation MUST be discovered, and instead places that burden upon the SP.
* Misconfigured plugin software SHOULD fail-fast with an OS-appropriate error code.


##### Namespaces

* A Plugin SHOULD NOT assume that it is in the same [Linux namespaces](https://en.wikipedia.org/wiki/Linux_namespaces) as the Provisioner Sidecar.
  The Database system MUST clearly document the CSI dependency requirements on ephemeral volumes for the plugins and the Provisioner Sidecar SHALL satisfy the Database system’s requirements.

##### Cgroup Isolation

* A Plugin MAY be constrained by cgroups.
* An operator or Provisioner Sidecar MAY configure the devices cgroup subsystem to ensure that a Plugin MAY access requisite devices.
* A Provisioner Sidecar MAY define resource limits for a Plugin.

##### Resource Requirements

* SPs SHOULD unambiguously document all of a Plugin’s resource requirements.
