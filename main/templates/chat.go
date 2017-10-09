package templates

import (
	"net/http"
	"fmt"
)

func Chat(w http.ResponseWriter, r *http.Request, user string) {

	fmt.Fprintf(w, `
<html>
<head>
    <title>%s Profile WebChat</title>
	<link type="text/css" rel="stylesheet" href="style.css" />
	<script src="http://code.jquery.com/jquery-1.11.3.min.js"></script>
</head>
<body>
	<center>
	<div id="wrapper">
    <div id="menu">
        <p class="welcome">Welcome, %s<b></b></p>
        <p class="logout"><a id="exit" href="/chat/logout">Exit Chat</a></p>
        <div style="clear:both"></div>
    </div>

    <div id="chatbox"></div>

    <input name="usermsg" type="text" autocomplete="off" id="usermsg" size="128" />
    <input type="button" id="submitmsg" value="Send" />
	</div>
	<center>

	<script>
    //Set console
    if(typeof window.console == 'undefined') { window.console = {log: function (msg) {} }; }

    // Get Recent Events
    var yourActionsSinceTime = (new Date(Date.now() - 120000)).getTime();
    var yourActionsCategory = "%s_actions";
    (function pollYourActions() {
        var timeout = 45;  // in seconds
        var optionalSince = "";
        if (yourActionsSinceTime) {
            optionalSince = "&since_time=" + yourActionsSinceTime;
        }
        var pollUrl = "api/requestMessage?timeout=" + timeout + "&category=" + yourActionsCategory + optionalSince;
        // how long to wait before starting next longpoll request in each case:
        var successDelay = 10;  // 10 ms
        var errorDelay = 3000;  // 3 sec
        $.ajax({ url: pollUrl,
            success: function(data) {
                if (data && data.events && data.events.length > 0) {
                    // got events, process them
                    // NOTE: these events are in chronological order (oldest first)
                    for (var i = 0; i < data.events.length; i++) {
                        // Display event
                        var event = data.events[i];
                        // prepend instead of append so newest is up top--easier to see with no scrolling
                        var message = "<div class='msgln'>(" + (new Date(event.timestamp).toLocaleTimeString()) + ") <b>" + event.data.user + "</b>: " + event.data.text + "<br></div>";
    					if (event.data.action == "whisper") {
							message = "<strong><div style=\"color:SlateBlue;\">" + message + "</div></strong>";
						}
						if (event.data.action == "notice") {
							message = "<strong>" + message + "</strong>";
						}
						$("#chatbox").append(message);
						$('#chatbox').animate({scrollTop: $('#chatbox').prop("scrollHeight")}, 500);
                        // Update sinceTime to only request events that occurred after this one.
                        yourActionsSinceTime = event.timestamp;
                    }
                    // success!  start next longpoll
                    setTimeout(pollYourActions, successDelay);
                    return;
                }
                if (data && data.timeout) {
                    console.log("No events, checking again.");
                    // no events within timeout window, start another longpoll:
                    setTimeout(pollYourActions, successDelay);
                    return;
                }
                if (data && data.error) {
                    console.log("Error response: " + data.error);
                    console.log("Trying again shortly...")
                    setTimeout(pollYourActions, errorDelay);
                    return;
                }
                // We should have gotten one of the above 3 cases:
                // either nonempty event data, a timeout, or an error.
                console.log("Didn't get expected event data, try again shortly...");
                setTimeout(pollYourActions, errorDelay);
            }, dataType: "json",
        error: function (data) {
            console.log("Error in ajax request--trying again shortly...");
            setTimeout(pollYourActions, errorDelay);  // 3s
        }
        });
    })();
    // Add another longpoller for all user's public events:
    var publicActionsSinceTime = (new Date(Date.now() - 120000)).getTime();
    var publicActionsCategory = "public_actions";
    // Longpoll subscription for everyone's (public) actions.
    // You wont see other people's private actions
    (function pollPublicActions() {
        var timeout = 45;  // in seconds
        var optionalSince = "";
        if (publicActionsSinceTime) {
            optionalSince = "&since_time=" + publicActionsSinceTime;
        }
        var pollUrl = "api/requestMessage?timeout=" + timeout + "&category=" + publicActionsCategory + optionalSince;
        // how long to wait before starting next longpoll request in each case:
        var successDelay = 10;  // 10 ms
        var errorDelay = 3000;  // 3 sec
        $.ajax({ url: pollUrl,
            success: function(data) {
                if (data && data.events && data.events.length > 0) {
                    // got events, process them
                    // NOTE: these events are in chronological order (oldest first)
                    for (var i = 0; i < data.events.length; i++) {
                        // Display event
                        var event = data.events[i];
                        var message = "<div class='msgln'>(" + (new Date(event.timestamp).toLocaleTimeString()) + ") <b>" + event.data.user + "</b>: " + event.data.text + "<br></div>";
                        if (event.data.action == "shout") {
							message = "<strong><div style=\"color:Orange;\">" + message + "</div></strong>";
						}
						if (event.data.action == "notice") {
							message = "<strong>" + message + "</strong>";
						}
						$("#chatbox").append(message);
						$('#chatbox').animate({scrollTop: $('#chatbox').prop("scrollHeight")}, 500);
                        // Update sinceTime to only request events that occurred after this one.
                        publicActionsSinceTime = event.timestamp;
                    }
                    // success!  start next longpoll
                    setTimeout(pollPublicActions, successDelay);
                    return;
                }
                if (data && data.timeout) {
                    console.log("No events, checking again.");
                    // no events within timeout window, start another longpoll:
                    setTimeout(pollPublicActions, successDelay);
                    return;
                }
                if (data && data.error) {
                    console.log("Error response: " + data.error);
                    console.log("Trying again shortly...");
                    setTimeout(pollPublicActions, errorDelay);
                    return;
                }
                // We should have gotten one of the above 3 cases:
                // either nonempty event data, a timeout, or an error.
                console.log("Didn't get expected event data, try again shortly...");
                setTimeout(pollPublicActions, errorDelay);
            }, dataType: "json",
        error: function (data) {
            console.log("Error in ajax request--trying again shortly...");
            setTimeout(pollPublicActions, errorDelay);  // 3s
        }
        });
    })();
//On enter make click
$("#usermsg").keyup(function(event){
    if(event.keyCode == 13){
        $("#submitmsg").click();
    }
});

// Click handler for action button.  This hits the http handler that publishes
// events.
$( "#submitmsg" ).click(function() {
    var action = "say";
    var message = document.getElementById("usermsg").value;
	document.getElementById("usermsg").value = "";
	var dest = "";

	//Shout
    if (message.startsWith("/shout")) {
        message = message.replace("/shout ", "");
        action = "shout";
    }

	//Exit
    if (message.startsWith("/exit")) {
        message = message.replace("/exit ", "");
        action = "exit";
    }

	//Promote
    if (message.startsWith("/promote")) {
        message = message.replace("/promote ", "");
		dest = message.split(' ')[0];
        message = message.replace(dest, "");
        action = "promote";
    }

	//Kick
    if (message.startsWith("/kick")) {
        message = message.replace("/kick ", "");
		dest = message.split(' ')[0];
        message = message.replace(dest, "");
        action = "kick";
    }

	//Ban
    if (message.startsWith("/ban")) {
        message = message.replace("/ban ", "");
		dest = message.split(' ')[0];
        message = message.replace(dest, "");
        action = "ban";
    }

	//Help
    if (message.startsWith("/help")) {
        message = message.replace("/help ", "");
        action = "help";
    }

	//Whisper
    if (message.startsWith("/w")) {
        message = message.replace("/w ", "");
        dest = message.split(' ')[0];
		action = "whisper";
        message = message.replace(dest, "");
    }

    var actionSubmitUrl = "api/sendMessage?user=" + dest + "&action=" + action + "&text=" + message;
	$.ajax({ url: actionSubmitUrl,
            success: function(data) {
            	console.log("action submitted");
            }, dataType: "html",
        error: function (data) {
        	//alert("Action failed due to error.");
        }
        });
});
</script>

</body>
</html>`, user, user, user)

}
