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
    - the auth token is to be sent with the header {"X-Auth-Token: <Session-Token>"}
## /new_post
### JSON query format:
	- Author string
	- Title string
	- Content string 
	- Tags string 
### HTTP Raw response format:
    - on succesful post creation will be: "Success" and code 200
    - on failure will be: <error that occured> and appropriate error code
## /logout
    just query it with aproprieate headers
## /submit_like
    - PostID int32 
    - LikerID int32
## /delete_like
    - PostID int32 
    - LikerID int32
