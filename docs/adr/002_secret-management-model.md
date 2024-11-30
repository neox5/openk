# ADR-002: Secret Management Model

## Status
Proposed

## Context
We need to define how secrets are organized, identified, and versioned in the system. This includes how users interact with secrets, how they are structured, and how versioning works.

## Decision

### 1. Secret Organization

#### 1.1 Path-based Hierarchy
* **Two Distinct Representations**:
  1. Internal/Visual Format:
     - Uses forward slash (`/`) for path separation
     - Leading and trailing slashes are removed during normalization
     - Used in UI, API responses, and internal processing
     - Example: `project/app`

  2. Reference String Format:
     - Uses underscore (`_`) for path segments
     - Used for CLI operations and technical references
     - Combined with labels and name in format: `<path>:<labels>:<name>`
     - Example: `project_app:env=prod,region=us-east-1:postgres-main`
     - Automatically converted to internal format after parsing

* **Character Constraints**:
  - Path segments: `[a-zA-Z0-9-_]+`
  - Label keys: Must match pattern `[a-zA-Z0-9][a-zA-Z0-9-_]*[a-zA-Z0-9]` or `[a-zA-Z0-9]`
                Maximum length: 63 characters
  - Label values: Must match pattern `[a-zA-Z0-9][a-zA-Z0-9-_.]*[a-zA-Z0-9]` or `[a-zA-Z0-9]`
                 Maximum length: 63 characters
  - Secret name: `[a-zA-Z0-9-_.]+`

* **Labels**: 
  - Identity qualifiers that are part of uniqueness
  - Some categories (like environments) might be predefined
  - Others can be user-defined (e.g., regions, datacenters)
  - Keys and values must follow Kubernetes-style constraints
  - Examples:
    ```
    Valid labels:
      env=prod
      region=us-east-1
      version=1.2.3
      app=my-service
    
    Invalid labels:
      .env=prod     (key starts with dot)
      env-=prod     (key ends with dash)
      env=prod.     (value ends with dot)
      env=-prod     (value starts with dash)
    ```

* **Tags**: 
  - Optional categorization helpers (e.g., pg, db)
  - Do not affect uniqueness

* **Uniqueness Rule**: 
  - Combination of `path + labels + name` must be unique
  - Label order does not affect uniqueness

#### 1.2 Path Examples

```plaintext
# Internal/Visual Format:
Path: project/backend
Labels: {"env": "prod"}
Name: database
Tags: [pg, db]

# Reference String Format:
project_backend:env=prod:database

# Path Normalization Examples:
Input: "/project/app/"    -> "project/app"
Input: "project/app/"     -> "project/app"
Input: "/project/app"     -> "project/app"
Input: "project/app"      -> "project/app"
```

### 2. Versioning

#### 2.1 Version Management
* Every modification creates a new version
* Latest version is always the active version
* Full versions stored for rollback capability
* Complete audit trail maintained

#### 2.2 Rollback Types
1. Hard Rollback:
   - Destroy newer versions until selected version
   - Cannot recover destroyed versions

2. Soft Rollback:
   - Copy previous version as new current version
   - Maintains full version history

## Consequences

### Positive
* Clear organizational structure with consistent path handling
* Distinct path notations for different use cases:
  - Forward slashes for human readability and UI
  - Underscores for technical operations
* Well-defined normalization rules
* Kubernetes-aligned label constraints for familiarity and proven patterns
* Full version history enables complete auditability
* Multiple rollback options support different use cases
* Character constraints ensure system-wide compatibility
* Support for both human and automated access patterns

### Negative
* Storage overhead from full version history
* Need to maintain path format conversion logic
* Potential confusion from multiple path formats
* Stricter label constraints may require updating existing labels

## Notes
* Label categories (like environments) may be predefined or user-definable
* Path conversion functions must be consistent across the system
* Need to document best practices for path structure in both notations
* Consider caching normalized paths for performance
* Future consideration: path-based access control