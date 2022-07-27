# message-api-sql
## Chatting Restfull API (like WhatsApp) GOlang with Go-Fiber   
   
### How To Run The App
1. Go to directory of the App (make sure u can run Go program)   
2. Then run this command
```
go run main.go
```   
3. It will looks like below picture if the program successfully running   
![cmd](https://user-images.githubusercontent.com/83685852/181220493-6e95a25c-9cc4-4922-b588-d89b821f5dea.png)   
   
### Case 1.1 : Authenticated User want to start chat and send message (with define parameter from form in the picture)
![story1](https://user-images.githubusercontent.com/83685852/181214276-1506155b-0d36-4431-a6a9-10423055f76a.png)  
### Case 1.2 : Authenticated User want to start chat and send message but the message empty (doesn't contain any char)   
![story1-1](https://user-images.githubusercontent.com/83685852/181217247-2b1b9c17-ac8d-412b-99d5-a14ac9820907.png)   
#### Form Value
1. id_user1 : id from authenticated user who want to start a chat message (e.g. user1) (required)
2. auth_code  : auth code/ credential from authenticated user (required)
3. message  : message value which had to be send, must contain at least one character (required)
4. id_user2 : id from destination user      
   
### Case 2 : Authenticated User want to respond message from other user (with define parameter from form in the picture)      
![story2](https://user-images.githubusercontent.com/83685852/181218268-1b7ed10b-90e0-4cec-bbfb-fbec1594872e.png)   
#### Form Value
1. id_user1 : id from authenticated user who want to reply to message (required)
2. auth_code  : auth code/ credential from authenticated user (required)
3. message  : message value which had to be send, must contain at least one character (required)
4. id_conv  : id from destination user      
   
### Case 3 : User want to listing all message from spesific conversation   
![story3](https://user-images.githubusercontent.com/83685852/181219040-cdbe0f7e-6627-491c-b227-ef0ba108ce32.png)   
#### URL/{id}
1. id : id from selected conversation   
   
### Case 4 : User want to listing all conversation they had   
![story4](https://user-images.githubusercontent.com/83685852/181219459-3ceea982-5574-4040-ab01-da28f7165b8a.png)      
#### URL/{id}
1. id : id from authenticated user      
