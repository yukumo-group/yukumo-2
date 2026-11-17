---
description: Sentinel errors re-exported from domain
icon: triangle-exclamation
---

# Errors

```go
var (
    ErrInvalid  = domain.ErrInvalid
    ErrNotFound = domain.ErrNotFound
    ErrExists   = domain.ErrExists
)
```

| Sentinel | Typical use |
| --- | --- |
| `ErrInvalid` | Bad input, missing engine catalog, invalid definition |
| `ErrNotFound` | Missing character/voice/config entry |
| `ErrExists` | Duplicate name on create |

{% hint style="info" %}
Prefer `errors.Is` against these sentinels rather than string matching.
{% endhint %}
