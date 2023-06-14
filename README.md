# assignment_demo_2023

<h2>This is my change to the demo provided to us</h2>
<h3>Docker Compose</h3>
Used MySQL as the database the RPC Server will connect to
<h3>The structure of the Database is<br/></h3>
CREATE TABLE Chatroom (<br/>
	id varchar(255),<br/>
    sender_id varchar(127),<br/>
    message varchar(255),<br/>
    time_stamp BIGINT,<br/>
    PRIMARY KEY (id, time_stamp)<br/>
)
<h3>HTTP-Server</h3>
Passed sent-time in microseconds to the RPC Server<br/>
Validated the chat-id to ensure that a:b and b:a refers to the same chatroom<br/>
Ensured that the chat-id is valid
<h3>RPC-Server</h3>
Handler.go contains all the functionality of the RPC Server
<h4>Handler.go</h4>
<h5>saveMessage</h5>
saveMessage contains the insertion query to store new messages
<h5>pullMessage</h5>
pullMessage contains the select query to retrieve messages<br/>
After retrieval, the messages are stored in the PullResponse and the PullResponse is updated depending on the overall status of the table