{{define "title"}}
Removed
{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}">
	Your registration has been removed if it was present.
	Our <strong>logs</strong> may still have information about you such as:
	<ul>
		<li>your IP address</li>
		<li>control codes (which cannot be used to identify your account on the Metropolis website without admin permissions)</li>
	</ul>
	Our <strong>database</strong> will not have any information about you, including your control code and OAuth tokens.
</div>
{{end}}
