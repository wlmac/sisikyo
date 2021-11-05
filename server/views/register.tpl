{{define "title"}}
Register
{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}">
	<h2>Instructions</h2>
	<ol>
		<li>press <a href="#authorize">"Authorize"</a></li>
		<li>review the permissions that you will grant to this app</li>
		<li>press "Authorize" on the next page to register</li>
	</ol>
	<a class="w3-button {{$card}}" id="authorize" href="{{.url}}">Authorize</a>
</div>
{{end}}
