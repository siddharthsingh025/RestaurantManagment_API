# RestaurantManagment_API

- To run this simply run

      go run main.go 

# All Routes :

      ## Menu 
            /menus
            /menus/:menu_id
      ## Order
            /orders
            /orders/:order_id
      ## Food
            /foods
            /foods/:food_id
      ## Table
            /tables
            /tables/:table_id
      ## OrderItem
            /orderItems
            /orderItems/:orderItem_id
            /orderItems-order/:order_id
      ## Invoice
            /invoices
            invoices/:invoice_id
      ## User
            /users
            /users/:user_id
            /users/signup
            /users/login

#### tips😄 : use genrated token to requests

# DataBase Schema I Used 

      [link]( )


[link]()


# Project Structure
    :RestaurantManagment 
      - Controllers         //controllers to do actions on datamodels and request
      - Database            //database connection 
      - Routes              //different routes according to request
      - Models              //contains DataModels 
      - Middleware          //Authentication service
      - Helper              //token genration service
      - go.mod              //go module
      - go.main             // server implimentation 

# Summary

    - RESTFul (approach/architechture) Microservice for  restaurantManagment
    - written in Golang using gin frameWork
    - DataBase used : MongoDb - noSQL database 
    - Authentication :JWT
    - Testing : PostMan , TablePlus
    - Deployment : docker compose , can be deployed in any cloud providers

# Testing
      
      ## using docker compose  ⭐
      
       1. docker compose up -d mongo-db     //running mongo-db with port 27017
       2. docker ps -a                     //check weather container is running or not
       3. docker compose build             // building our go-app - it will use our Dockerfile for it
       4. docker compose up go-app         //running our restaurantManagment server at post 8000 {go-app}
       5. docker ps -a                     // you can check weather both services are running or not as container


 ## docker compose building output  

<img width="974" alt="Screenshot 2023-07-17 at 2 26 08 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/26618853-4c7b-48a6-a111-9696279242e1">

<img width="991" alt="Screenshot 2023-07-17 at 2 27 00 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/d769171b-9702-4b38-825d-503d052ce8bd">

<img width="1014" alt="Screenshot 2023-07-17 at 2 27 19 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/fbb11118-e4db-4854-aac6-1c6e9a3dd0fb">

<img width="1220" alt="Screenshot 2023-07-17 at 2 28 17 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/eaa4e307-7810-4738-9b90-5aa37c50497a">

## table plus output 

<img width="1440" alt="Screenshot 2023-07-17 at 1 22 27 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/24fbb344-5a34-47ec-b9d7-b867befb9e8b">

<img width="1440" alt="Screenshot 2023-07-17 at 1 59 10 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/cc601bdb-29d8-42db-b065-3e632bfa1003">
<img width="1440" alt="Screenshot 2023-07-17 at 1 59 21 PM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/30690650-b57b-4317-af56-ab8caa27ba05">
     
 <img width="367" alt="Screenshot 2023-07-17 at 1 09 59 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/5d585372-3408-4b2d-8474-cf63fc3c193b">

# OutPut ScreenShorts :
<img width="1440" alt="Screenshot 2023-07-17 at 1 35 19 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/f3d33077-2273-4cd7-b860-01b18999c88c">
<img width="1440" alt="Screenshot 2023-07-17 at 1 14 05 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/3ada7c24-aaa4-4e74-bf29-9bdc696d1782">
<img width="1440" alt="Screenshot 2023-07-17 at 1 44 44 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/978d1bb2-b35f-427e-a0f0-079c971b71a1">
<img width="1440" alt="Screenshot 2023-07-17 at 1 44 19 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/4592d551-f17e-4ce8-b7b2-d0cf08192115">
<img width="1440" alt="Screenshot 2023-07-17 at 1 42 58 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/2d573aff-868d-461f-a93c-9c59ce5bd426">
<img width="1067" alt="Screenshot 2023-07-17 at 1 28 50 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/6e850849-c70d-412b-bb60-79f6f34cdaa1">
<img width="1226" alt="Screenshot 2023-07-17 at 2 00 52 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/ef81c1c6-456f-414a-bd69-8cdb1406d4e8">
<img width="1238" alt="Screenshot 2023-07-17 at 2 00 45 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/6041c281-c96e-40aa-888f-a64f07da16c6">
<img width="1229" alt="Screenshot 2023-07-17 at 2 00 38 AM" src="https://github.com/siddharthsingh025/RestaurantManagment_API/assets/87073574/bdcbe711-5ce8-42c5-b2a5-c1aa4e2c779e">





















