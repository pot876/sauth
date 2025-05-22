# sauth (simplified authorization/sessions logic)

#### endpoints:
```
POST /login
POST /refresh

todo POST /logout

todo POST /validate
todo GET  /info
```

#### Login:
```
--> { "realm_id": "${DOMAIN_ID}", "login": "", "password": "" }
<-- { access_token: "", refresh_token: "" }
```

#### Refresh:
```
--> { "token": "${REFRESH_TOKEN}" }
<-- { access_token: "", refresh_token: "" }
```

### Features TODO:
- Users sessions invalidation and analitics
- Users info invalidation
- Token leaks analisis, suspicious activity analisis
- DDOS protection
- Seamless keys rotation
- Prevents timing brute-force attack


---

#### generate keys
```sh
ssh-keygen -t rsa -b 2048 -m PEM -f jwtRS256.key -N ""
openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub

echo "export AUTH_ACCESS_TOKEN_KEY_ID=2aec25bd-ced1-4de9-9b2f-ff33a03b21ed"
echo "export AUTH_ACCESS_TOKEN_PRIVATE_KEY=$(cat jwtRS256.key | base64)"
echo "export AUTH_ACCESS_TOKEN_PUBLIC_KEY=$(cat jwtRS256.key.pub | base64)"

echo "export AUTH_REFRESH_TOKEN_KEY_ID=68758643-62f0-43bd-ab67-f0183c8a7eab"
echo "export AUTH_REFRESH_TOKEN_PRIVATE_KEY=$(cat jwtRS256.key | base64)"
echo "export AUTH_REFRESH_TOKEN_PUBLIC_KEY=$(cat jwtRS256.key.pub | base64)"
```


#### example
```sh
HTTP_LISTEN_ADDR=:8085 \
CP_FILE_ENABLED=true \
CP_FILE_PATH=internal/delegates/credentials_provider_file/example.json \
    go run ./cmd/auth_server


curl -i -X POST localhost:8085/login -d \
    '{"login":"alice", "password": "test", "realm_id": "3d038b0f-bcde-443a-817d-68f6723699b9"}'
```
