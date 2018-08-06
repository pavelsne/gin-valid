package templates

var Login = `
{{ define "content" }}

Please log in using your GIN credentials
<hr>
<div>
    <form action="/login" method="post">
        <p>Username: <input type="text" name="username"></p>
        <p>Password: <input type="password" name="password"></p>
        <input type="submit" value="Submit">
    </form>
</div>


{{ end }}
`
