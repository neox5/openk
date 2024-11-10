# StorageBackend

```go
type StorageBackend interface {
    CreateSecret(secret Secret) (string, error)
    GetSecret(id string, version int) (Secret, error)
    UpdateSecret(id string, newSecret Secret) error
    DeleteSecret(id string) error
}
```
