package templates

import (
	"net/http"
	"fmt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, `
<html>
<head>
    <title>WebChat Login Page</title>
</head>
<body>
	<center>
	<h2>Wait while login...</h2>

	<script>
	window.location.href = '/chat/main';
	</script>

</body>
</html>`)

}