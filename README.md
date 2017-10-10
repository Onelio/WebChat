# WebChat
WebChat is a Server made in Go using long poll. 

This project has been developed using jcuga golongpoll library which can be found here: https://github.com/jcuga/golongpoll


As default user connected with username Onelio is going to be admin. Remember to change it!

# Usefull Commands
## Everybody
/help To get all commands (In development)

/w <nick> <message> To whisper someone so no one else can see it
## Admin
/shout <message> To say something resalted(in yellow) from the rest of the users

/promote <nick> To promote a user as admin(Sending it again makes a admin user)

/kick <nick> To kick an user from the chat

/ban <nick> To ban an user from the chat (Still in development)

/exit To close the chat-server
