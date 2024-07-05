# For Radicale server:
add user:
`htpasswd ./radicale/users someusername`


start Radicale:
`./runRadicale.sh`


build and run client:
`make run`


before running tests remember to start Radicale

# For Baikal server:
start Baikal:
`docker compose up`


add user:
`go to 127.0.0.1:90, register as admin, add users`


build and run client:
`make run`


before running tests make sure to edit test user data in caldav_test.go
