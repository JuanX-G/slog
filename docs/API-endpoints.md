# Unrestricted endpoints

## /get_user_posts
### JSON query format:
    -   UserName string 
    -   Count int 
    -   Offset int 

### JSON response format:
    - Content string
    - Title string
    - DatePosted time.Time
    - Tags string
    - Likes int32
    - ID int32 

## /new_user
### JSON query format:
    - UserName string 
    - Email string
    - DateCreated time.Time 
    - Password string
    - Description string
### JSON response format:
    - Response string 
        Will be: "Account created succesfully" upon succesfull account creation

## /login
### HTTP Form request format:
        key         |       Value       
       user         :     <username>
       pass         :     <password>
### HTTP Raw response format:
    <token-for-the-session>

## /get_user_description
### JSON query format:
    - UserName string
### JSON response format:
    -  Description string


# Endpoints requiring auth
/new_post
/logout
/submit_like
/delete_like
