# otus_architect_social

## Local running
*Run main project*:
```make up```

With default settings web-interface will be on http://localhost:8080/
Server service starts on :8090

*Switch off*:
```make down```

*Run messenger server*:
```make up-messenger```

Server starts on local 8091 port.

*Switch off messenger server*:
```make down-messenger```

### Examples (terminal)
1. Register user:

```curl -d '{"email": "test@email.com", "password": "...", "first_name": "TestFirst", "last_name": "TestLast", "birth_date": "1990-11-12T00:00:00Z", "sex": "M", "interest": "list of interests", "city": "Moscow"}' -H "Content-Type: application/json" -X POST http://:8090/api/register```

Server returns access-token for this request. You have to use it in other requests that require authentification. You need send it through header ```Authorization: Bearer [access-token]```.

2. Send message to user_id=90002:

```curl -d '{"to_user_id": 90002, "body": "test message body3"}' -H "Content-Type: application/json" -H "Authorization: Bearer ..." -X POST http://:8091/api/send```

3. Get dialog messages for user_id=90002:

```curl -H "Content-Type: application/json" -H "Authorization: Bearer ..." -X GET http://:8091/api/dialog/get/90002```
