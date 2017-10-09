package templates

import (
	"net/http"
	"fmt"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<html>
<head>
    <title>WebChat by Onelio</title>
</head>
<body>
    <h1>WebChat Server made with Go by Onelio</h1>
    This is a webchat that implements a server made in Go and uses long poll to handle messages.
	<br>
	<br>
	<center>
	<h4>Please set a nickname to continue:</h4>
	<input type="text" id="nick" value="Unknown"> <button onclick="login()" type="button">Connect</button>

	<script>
	function login() {
		var user = document.getElementById("nick").value;
		window.location.href = '/chat?user=' + user;
	}
	</script>

</body>
</html>`)
}