# Chatbox API 👋

### User Registration

```
HTTP Method: POST
URL: {{url}}/api/v1/auth/
```

##### Sample Request Body

```
{
    "email": "user1@example.com",
    "password": "12345678",
    "password_confirmation": "12345678"
}
```

##### Parameters

| Name                  | Description      | Required |
| --------------------- | ---------------- | -------- |
| email                 | User email       | Yes      |
| password              | User password    | Yes      |
| password_confirmation | Re-type password | Yes      |

### Login

```
HTTP Method: POST
URL: {{url}}/api/v1/auth/sign_in
```

##### Sample Request Body

```
{
    "email": "meline@hotmail.com",
    "password": "12345678"
}
```

##### Parameters

| Name     | Description   | Required |
| -------- | ------------- | -------- |
| email    | User email    | Yes      |
| password | User password | Yes      |

### Send Message

```
HTTP Method: POST
URL: {{url}}/api/v1/messages
```

##### Sample Request Body

```
{
    "receiver_id": 1,
    "receiver_class": "User",
    "body": "kamusta?"
}
```

##### Parameters

| Name           | Description                                                                                   | Required |
| -------------- | --------------------------------------------------------------------------------------------- | -------- |
| receiver_id    | ID of the message's receiver                                                                  | Yes      |
| receiver_class | Type of the receiver. `User` for direct message, `Channel` for sending a message in a channel | Yes      |
| body           | Message body                                                                                  | Yes      |

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### Retrieve Message

```
HTTP Method: Get
URL: {{url}}/api/v1/messages?receiver_id=1&receiver_class=User
```

##### Parameters

| Name           | Description                                                                                   | Required |
| -------------- | --------------------------------------------------------------------------------------------- | -------- |
| receiver_id    | ID of the message's receiver                                                                  | Yes      |
| receiver_class | Type of the receiver. `User` for direct message, `Channel` for sending a message in a channel | Yes      |

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### Create Channel with members

```
HTTP Method: POST
URL: {{url}}/api/v1/channels
```

##### Sample Request Body

```
{
    "name": "channel1",
    "user_ids": [2]
}
```

##### Parameters

| Name     | Description                                           | Required |
| -------- | ----------------------------------------------------- | -------- |
| name     | Channel Name                                          | Yes      |
| user_ids | List of user ids to be included on the channel. Array | Yes      |

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### Get all users channels

```
HTTP Method: Get
URL: {{url}}/api/v1/channels
```

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### Get channel details via channel ID

```
HTTP Method: Get
URL: {{url}}/api/v1/channels/3
```

##### Parameters

| Name | Description       | Required |
| ---- | ----------------- | -------- |
| id   | ID of the channel | Yes      |

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### Add member to a channel

```
HTTP Method: POST
URL: {{url}}/api/v1/channel/add_member
```

##### Sample Request Body

```
{
    "id": 3,
    "member_id": 3
}
```

##### Parameters

| Name      | Description               | Required |
| --------- | ------------------------- | -------- |
| id        | Channel ID                | Yes      |
| member_id | User ID of the new member | Yes      |

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |

### List of All Users

```
HTTP Method: Get
URL: {{url}}/api/v1/users
```

##### Request Headers

_Get these values from the login response header_

| Name         | Required |
| ------------ | -------- |
| access-token | Yes      |
| client       | Yes      |
| expiry       | Yes      |
| uid          | Yes      |
