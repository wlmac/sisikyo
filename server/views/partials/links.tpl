{{define "links"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}} {{.class}}">
	<h2>Links</h2>
	<h3>Public</h3>
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/public.ics">Download as iCalendar file</a>
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/public.html">View</a>
	{{/*
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/public.json">Download as JSON file</a>
	*/}}
	{{with .control}}
	<h3>Private</h3>
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/{{.}}/private.ics">Events (vCalendar)</a>
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/{{.}}/private.html">Events (HTML)</a>
	{{/*
	<a class="w3-margin w3-padding w3-round w3-button w3-border" href="/events/{{.}}/private.json">Events (JSON)</a>
	*/}}
	{{end}}
</div>
{{end}}
