{
    "firstname": "test",
    "lastname": "user",
    "username": "user12",
    "email":    "user@gmail.com",
    "password": "Password123",
    "confirm-password": "Password123"
}



curl -X POST \
  http://localhost:3000/register \
  -H 'Content-Type: application/json' \
  -d '{
    "firstname": "test",
    "lastname": "user",
    "username": "user12",
    "email":    "user@gmail.com",
    "password": "Password123",
    "confirm-password": "Password123"
}'