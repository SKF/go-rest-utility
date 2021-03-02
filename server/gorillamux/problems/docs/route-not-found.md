[< Index](/problems)

# The requested endpoint could not be found

The request was unable to be routed to one of the specified endpoints.

## Fields

|           |                                                                                       |
| --------- | ------------------------------------------------------------------------------------- |
| `type`    | Path to this page                                                                     |
| `title`   | The title of this Problem                                                             |
| `status`  | HTTP status code that is returned with this Problem                                   |
| `detail`  | Message about why this Problem was returned                                           |

## Example

```json
{
    "type": "/problems/route-not-found",
    "title": "The requested endpoint could not be found.",
    "status": 404,
    "detail": "Ensure that the URI is a valid endpoint for the service."
}
```
