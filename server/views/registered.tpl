{{define "head"}}Registered{{end}}
{{define "body"}}
{{$card := "w3-card w3-margin w3-padding w3-round"}}
<div class="{{$card}}">
	<h2>Add Events to Google Calendar</h2>
	Follow
	<a href="https://support.google.com/calendar/answer/37100?hl=en#:~:text=Use%20a%20link%20to%20add%20a%20public%20calendar">
		this guide about how to add the link below
	</a>
	to your calendar.
	<ul>
		<li><a href="/events/{{.control}}/private.ics">Link for Private Events</a></li>
		<li><a href="/events/public.ics">Link for Public Events</a></li>
	</ul>
</div>
{{template "links" .}}
{{end}}
