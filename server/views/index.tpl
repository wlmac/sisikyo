{{define "title"}}
Sisikyō
{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
{{template "links" .}}
<div class="{{$card}}" id="debug">
	<h2>Debug Info</h2>
	API Base URL: <pre>{{.apiURL}}</pre>
	OAuth Base URL: <pre>{{.oauthBaseURL}}</pre>

	{{with .buildInfo}}
	<h3>Build Info</h3>
	Path: <pre>{{.Path}}</pre>
	Main: {{.Main.Path}} @ {{.Main.Version}}{{with .Main.Sum}} (sum <pre>{{.}}</pre>){{end}}{{if ne .Main.Replace nil}}<span class="w3-tag w3-round">replaced</span>{{end}}
	<h4>Details</h4>
	<pre>{{$.index}}</pre>
	{{end}}
</div>
{{end}}
