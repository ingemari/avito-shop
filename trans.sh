TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3MzkzODM4ODN9.9JPsbnuT9-_mKggM84Wbx28XG_7nBaIhmRJeFCXrclc"

curl -X POST http://localhost:8080/api/transaction \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
           "to": "testuser",
           "amount": 100
         }'