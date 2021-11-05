{{define "links"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}">
	<h2>Links</h2>
	<h3>Public</h3>
	<a class="{{$card}} w3-button" href="/events/public.ics">Events (vCalendar)</a>
	<a class="{{$card}} w3-button" href="/events/public.json">Events (JSON)</a>
	<a class="{{$card}} w3-button" href="/events/public.html">Events (HTML)</a>
	{{with .control}}
	<h3>Private</h3>
	<a class="{{$card}} w3-button" href="/events/{{.}}/private.ics">Events (vCalendar)</a>
	<a class="{{$card}} w3-button" href="/events/{{.}}/private.json">Events (JSON)</a>
	<a class="{{$card}} w3-button" href="/events/{{.}}/private.html">Events (HTML)</a>
	{{end}}
</div>
{{end}}
