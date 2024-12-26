# gRPC Style Guide

## Buf Configuration

### buf.yaml
```yaml
version: v2
modules:
  - path: .
    name: buf.build/neox5/openk
    excludes:
      - vendor/google/protobuf
lint:
  use:
    - DEFAULT
  except:
    - PACKAGE_VERSION_SUFFIX
breaking:
  use:
    - FILE
```

Key settings:
- `modules`: Local module configuration without external dependencies
- `lint`: Enforces API design standards (DEFAULT ruleset)
- `breaking`: Detects breaking API changes using FILE ruleset
- `excludes`: Prevents vendor directory from being processed

### buf.gen.yaml
```yaml
version: v2
clean: true
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/neox5/openk/internal/api_gen
plugins:
  - remote: buf.build/protocolbuffers/go:v1.31.0
    out: ../internal/api_gen
    opt: 
      - paths=source_relative
  - remote: buf.build/grpc/go:v1.3.0
    out: ../internal/api_gen
    opt:
      - paths=source_relative
```

Key settings:
- `clean`: Clears output directory before generation
- `managed`: Enables consistent package naming
- `plugins`: Configures Go and gRPC code generation
- `out`: Places generated code in internal/api_gen

## Directory Structure
```
openk/
├── proto/                      # Proto definitions
│   ├── buf.yaml               
│   ├── buf.gen.yaml           
│   ├── openk/                 # API namespace
│   │   └── <service>/        # Service directories
│   │       ├── v1/           # Version directories
│   │       │   ├── service.proto  # Service definitions
│   │       │   └── types.proto    # Type definitions
│   │       └── v2/
│   └── vendor/               # Vendored dependencies
│       └── google/
│           └── protobuf/
└── internal/
    └── api_gen/             # Generated code
        └── openk/
            └── <service>/
                └── v{n}/
```

## Package Naming
- Format: `openk.<service>.<version>`
- Example: `openk.health.v1`
- Never reuse package names across different protofiles
- Always include version in package name

## Go Package Naming
- Format: `github.com/neox5/openk/internal/api_gen/openk/<service>/<version>;<service><version>`
- Example: `github.com/neox5/openk/internal/api_gen/openk/health/v1;healthv1`
- Import path matches generated code location
- Short package names for clean imports
