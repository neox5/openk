## Context & Logging Implementation Todos

### Service Information
- [ ] Implement service version from build info
  - Research Go build flags for version injection
  - Create build script/makefile target
  - Add version info package
- [ ] Instance ID implementation
  - Define instance ID format/requirements
  - Implement generation strategy
  - Add configuration options
  - Consider cloud platform integration

### Configuration
- [ ] Service information configuration
  - Move constants to config package
  - Add environment variable support
  - Add validation

### Main Implementation
- [ ] Review initialization order
  - Context vs Logger initialization
  - Error handling during startup
  - Graceful degradation options

### Future Considerations
- Context-aware logger wrapper
- Service discovery integration
- Cloud platform metadata integration
- Instance ID rotation/lifecycle