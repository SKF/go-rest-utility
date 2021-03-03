[< Index](/problems)

# The requested method is not allowed

The API endpoint did not accept the `HTTP` method that you requested. Verify that,

-   The endpoint you tried to call exists in the specification.
-   The `requestedMethod` in the problem matches the one in the specification.

## Fields

|                    |                                                                                       |
| ------------------ | ------------------------------------------------------------------------------------- |
| `requestedMethod`  | The requested method that is not allowed                                              |
| `allowedMethods`   | List of all accepted methods for this endpoint. Same as provided the `Accept` header. |

## Example

```json
{
    "type": "/problems/request-method-not-allowed",
    "title": "The requested method is not allowed.",
    "status": 405,
    "detail": "The requested resource does not support method 'PATCH', it does only support one of 'GET, POST, PUT, DELETE'.",
    "requestedMethod": "PATCH",
    "allowedMethods": ["GET", "POST", "PUT", "DELETE"]
}
```
