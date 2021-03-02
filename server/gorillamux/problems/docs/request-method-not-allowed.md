[< Index](/problems)

# The requested method is not allowed

The API endpoint did not accept the `HTTP` method that you requested. Verify that,

-   The endpoint you tried to call exists in the specification.
-   The `method` in the problem matches the one in the specification.

## Fields

|           |                                                                                       |
| --------- | ------------------------------------------------------------------------------------- |
| `method`  | The requested method that is not allowed                                              |
| `allowed` | List of all accepted methods for this endpoint. Same as provided the `Accept` header. |
| `type`    | Path to this page                                                                     |
| `title`   | The title of this Problem                                                             |
| `status`  | HTTP status code that is returned with this Problem                                   |
| `detail`  | Message about why this Problem was returned                                           |

## Example

```json
{
    "type": "/problems/request-method-not-allowed",
    "title": "The requested method is not allowed.",
    "status": 405,
    "detail": "The requested resource does not support method 'PATCH', it does only support one of 'GET, POST, PUT, DELETE'.",
    "method": "PATCH",
    "allowed": ["GET", "POST", "PUT", "DELETE"]
}
```
