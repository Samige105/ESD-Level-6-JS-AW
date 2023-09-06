# ESD Level 6 JS/AW
 
## Collaborators
###### James Sadler
###### Aidan Willis

## Tasks
* Create Repository - Done
* Create diagrams:
    * ERD for database
    * Other?    
@@ -0,0 +1,172 @@
# ITPR6.518 Assessment notes

You will build an application in the Go language, potentially in a team of 2 people. It is important you will clearly identify your contribution to the project. You will get an individual mark based on your contribution and your role in the team.

The application you will build is called Enterprise Notes and offers the notion of a shared notebook with a task tracking feature. The idea is that notes and tasks can be shared with identified users who can read and write notes in this application. You will pay special attention to a search and filter function that enables users to sift through unstructured texts. Security is an important feature as only authorised users on this notebook will have access. The application will be available as a web-service with a light-weight front-end (you're not required to build a comprehensive GUI front-end).


## ACTION PLAN

1. Creating the project folder
1. Establish a git repository
1. Creating the main.go, 
1. Writing the CRUD functions

## Tools & Libraries used

- lib/pq or jackc/pgx for the database
    - [demo](https://github.com/yonush/jobs.tradie.pg)
    - *todo* demo
- [optional] KV datastore Clover
    - [demo](https://github.com/yonush/jobs.tradie.kv)
- 

## Requirements

- A note just contains text (embedded media like video, images, etc is not required). 
- A note should have a name to aid with identification of the note.
- The note should contain a date & time of creation. 
- And include a date & time of completion if it is a task. 
- The note/task should include a status flag to indicate if the task is none/in progress/completed/cancelled/delegated.
- User id or name of the user the note was delegated to. (**ambiguous**)

### Required database tables

- Notes - store all of the notes regardless of the user/owner
    - ID
    - title
    - content
    - timestamp
    - status (enumerated type)
    - userID (owner of the note)

- Users - all of the uses of the software - owner + reader/guest
    - user ID
    - name (first+lastname)
    - email
    - password hash

- Sharing - which notes are shared with whom
    - user ID
    - note ID
    - status (enumerated type)
    - timestamp

Relation: One user with many Notes
Relation: One Note can be shared with many users


### 1. Find a note
    
- Refer to 2. below.
- Using a SELECT query

### 2. Analyse a note

- Using a SELECT query

### 3. Read/Edit a note

- Using a SELECT query
- Using a UPDATE query


## Architecture

- Low level CRUD (see section below)
- Driver code for the CRUD - can be done via _test files
- RESTful API - test with _test or curl scripts
    - gorilla/mux
    - refer to *Restful API* demos
- UI - us Go HTML templates
    - refer to the *jobs.tradie.kv* or *jobs.tradie.pg* & *famcost* demos     
    - [go templates](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go)
    - [tutorial 1](https://blog.logrocket.com/using-golang-templates/)
    - [tutorial 2](https://zetcode.com/golang/template/)

## CRUD

- Creating functions that can handle the data to/from the database.
- Includes user permissions where necessary

### Manage a note

Refer to the *todo* & the *jobs.tradie.pg* example

- Create a note
    ```go
        type Notes struct {
            id int
            title string
            content []string
            status string
            userID int
        }

        func (n *Notes) createNote(N Notes) {
            ...
        }
    
    ```

- Retrieve a note
- Edit a note
- Delete a note (**not mentioned in the brief**)
- Query a note

### Manage a user

(**not mentioned in the brief**)

- Create a user
- Retrieve a user
- Edit a user
- Delete a user

## Driver code for the CRUD

Use the stock *testing* package from Go

Explore the use of the following package for testing
 - [testify](https://github.com/stretchr/testify)

### Notes driver/test code

- Create a note
    ```go
        func TestNoteAdd(t *testing.T) {
	        SQL := ...
            expected := ...
            // Execute the SQL - add the necessary code
	        if observed != expected {
		        t.Errorf("Expected behaviour did not match observe behaviour!")
	        }
            
        }        
    
    ```

- Retrieve a note
- Edit a note
- Delete a note (**not mentioned in the brief**)
- Query a note

### Users driver/test code

## Deliverables    

Discuss features of the Enterprise Notes that are candidates to be executed on the client side instead of on the server. 

- Clearly describe the pros and cons.
- All persistent data (notes, user accounts, etc.) are to be stored in a suitable datastore. Explain the design 
choices you made to interact with the database.
- The enterprise typically hosts a variety of operating systems and internet browsers. Discuss how your 
solution copes with this variety.
- Provide a "Quick Start Guide" outlining the steps and details required to install your application on a new server.
- Provide a document that lists the additional specifications that were missing but required to implement 
your solution.




-end-of-file-