# OpenK ADR Writing Guide

## Purpose
This guide defines how to write clear and concise Architecture Decision Records (ADRs) for OpenK.

## ADR Structure

### Required Sections
```markdown
# ADR-[number]: [Title]

## Status
[Proposed, Accepted, Revised, Superseded]

## Context
[1-3 sentences describing the core problem and constraints]

## Decision
[Core technical decisions with clear structure]

## Consequences
[Clear positive and negative impacts]
```

### Optional Sections
- **References**: Only if directly referenced in decisions
- **Notes**: Very brief implementation guidance if needed
- **Migration**: Only for changes affecting existing systems

## Writing Guidelines

### Title
- Start with domain/category: "Crypto: Key Derivation"
- Use active voice: "Manage Keys" not "Key Management"
- Max 6 words

### Status
- Single word only: Proposed, Accepted, Revised, Superseded
- If Revised/Superseded, reference new/old ADR number
- Example: "Revised (see ADR-015)"

### Context 
- Max 3 sentences
- Structure: Problem + Constraints + Impact
- Focus on technical facts only
- No history or background

### Decision
- Structure content by topic with clear headings
- Use tables when presenting:
  - Configuration parameters
  - Constraints and their values
  - Comparison of options
- Include minimal code examples when they clarify
- Each decision must be actionable

### When to Use Tables
Tables are appropriate for:
- Sets of configuration parameters
- Mapping of inputs to outputs
- Comparison of multiple options
- Constraint definitions

Tables are not needed for:
- Simple key-value pairs
- Sequential steps
- Single option descriptions
- Implementation details

### Consequences
Present impacts in two clear categories:
- Positive impacts (benefits, advantages)
- Negative impacts (costs, limitations)

Focus on architectural and technical impacts, not process or organizational changes.

## What Not To Include
- Implementation details
- Obvious security practices
- Product comparisons
- Future considerations
- Background information
- Extended explanations
- Non-technical impacts

## Review Checklist
1. Title clearly identifies domain and decision
2. Context is exactly 1-3 sentences
3. Decisions are clearly structured
4. Tables used only where they add clarity
5. Code examples are minimal and clear
6. Consequences clearly separated into positive/negative
7. No banned content (see above)
8. Can be fully read in < 5 minutes

## Example - Well-Structured ADR

```markdown
# ADR-000: Crypto: Derive User Keys

## Status
Accepted

## Context
Need deterministic key derivation from user credentials that separates authentication from encryption operations while maintaining simplicity.

## Decision

### Key Derivation Parameters
| Parameter    | Value    | Constraints        |
|-------------|----------|--------------------|
| Algorithm   | PBKDF2   | FIPS 140-2         |
| Iterations  | 100,000  | Min 100K           |
| Salt        | Username | ASCII only         |
| Output      | 256 bits | For AES-256-GCM    |

### Implementation Strategy
Key generation follows a two-step process:
1. Generate master key from user credentials
2. Derive authentication key from master key

Core implementation:
```go
// Core derivation
masterKey := pbkdf2.Key(password, username, 100000) 
authKey := pbkdf2.Key(masterKey, AUTH_SALT, 100000)
```

## Consequences

Positive:
- No salt storage required
- Clear separation of authentication and encryption
- Minimal state tracking

Negative:
- Username changes affect derived keys
- Fixed to PBKDF2 algorithm
- Server-side computation required
```

## Common Mistakes
1. Including background/history in Context
2. Mixing decisions with explanations
3. Adding future considerations
4. Including obvious/standard practices
5. Writing long prose explanations
6. Adding implementation details
7. Including non-technical impacts
8. Overusing tables where simple lists suffice
